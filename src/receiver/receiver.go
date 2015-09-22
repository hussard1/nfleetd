package main

import(
//	"fmt"
	"runtime"
	"net"
	"sync"
	"regexp"
	"github.com/spf13/viper"
	"strconv"
	log "github.com/cihub/seelog"
	"os"
	"encoding/hex"
	"gopkg.in/mgo.v2"
	"time"

)

const (
	//config file name
	cofingfile = "deviceInfo"
	//config file path
	configpath = "\\resource\\config\\"
	//log config file name
	logconfigfile = "seelog.xml"
)


type DeviceInfo struct{
	name string
	status bool
	protocol string
	port int
	regex string
	threadCnt int
	buffer int
}

type Database struct{
	session *mgo.Session
}


func main() {
	//prevent to stop main goroutine
	done := make(chan struct{})
	defer close(done)

	//use all CPU core
	numCPUs := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPUs)

	//init logger
	initLogger()
	defer log.Flush()

	//init Database
	database := new(Database)
	database.session = InitMongoDB()

	//get config
	err := getDeviceInfoFile()

	if err != nil{
		return
	}else{
		for _, name := range viper.AllKeys() {
			deviceInfo := new(DeviceInfo)
			deviceInfo.initDevice(name)
			if deviceInfo.status == true {
				go func(){
					startReciever(deviceInfo, database)
				}()
			}
		}
	}

	for _ = range done{
		select{
			case <-done :
				return
		}
	}
}

func initLogger(){
	path, err := os.Getwd()
	if err != nil{
		return
	}
	logger, err := log.LoggerFromConfigAsFile(path + configpath + logconfigfile)
	if err != nil{
		logger.Info("Log Config file error : ", err)
		return
	}
	log.ReplaceLogger(logger)
	log.Info("Get config file : " + path + configpath + logconfigfile)
}

func getDeviceInfoFile() error{

	path, err := os.Getwd()
	if err != nil{
		log.Info("Config file error : ", err)
	}
	//read Configuration file
	viper.SetConfigName(cofingfile) // name of config file (without extension)
	viper.AddConfigPath(path + configpath)      // path to look for the config file in

	err = viper.ReadInConfig()
	if err != nil {
		log.Info("Config file error : ", err)
	}
	log.Info("Get config file : " + path + configpath + cofingfile)
	return err
}

func (deviceInfo *DeviceInfo) initDevice(name string){
	deviceInfo.name = name
	deviceInfo.status = viper.GetBool(name+".status")
	deviceInfo.protocol = viper.GetString(name+".protocol")
	deviceInfo.port = viper.GetInt(name+".port")
	deviceInfo.regex = viper.GetString(name+".regex")
	deviceInfo.threadCnt = viper.GetInt(name+".threadCnt")
	deviceInfo.buffer = viper.GetInt(name+".buffer")
}


func startReciever(deviceInfo *DeviceInfo, database *Database) {

	log.Info("start to receive data : name=" + deviceInfo.name + " port=" + strconv.Itoa(deviceInfo.port))

	ln, err := net.Listen(deviceInfo.protocol, ":" + strconv.Itoa(deviceInfo.port))
	if err != nil{
		log.Critical(err)
		return
	}
	defer ln.Close()

	ch := make(chan []byte)

	var wg sync.WaitGroup
	numWorkers := deviceInfo.threadCnt
	wg.Add(numWorkers)

	for i:=0; i< numWorkers; i++{
		go func(n int) {
			worker(n, ch, deviceInfo, database)
			wg.Done()
		}(i)
	}

	go func() {
		wg.Wait()
	}()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Info(err)
			continue
		}
		defer conn.Close()

		go func(c net.Conn) {
			data := make([]byte, 4096)
			for {
				n, err := c.Read(data)
				if err != nil {
					log.Critical(err)
					return
				}
				ch <- data[:n]
			}
		}(conn)
	}
	defer close(ch)
}

func worker(n int, ch<-chan []byte, deviceInfo *DeviceInfo, database *Database){
	pData := make(map[string]string)
	for rawData := range ch{
		select{
			default :
			pData = parseDataRegex(hex.EncodeToString(rawData), deviceInfo.regex);
			InsertMapToMongoDB(5, pData, database.session)
		}
	}
}

type myRegexp struct {
	*regexp.Regexp
}

func parseDataRegex(rawData string, regex string) map[string]string{
	s1 := make(map[string]string)
	re1 := myRegexp{regexp.MustCompile(regex)}
	s1 = re1.FindStringSubmatchMap(rawData)
	return s1
}


func (r *myRegexp) FindStringSubmatchMap(s string) map[string]string{
	captures := make(map[string]string)

	match := r.FindStringSubmatch(s)
	if match == nil{
		return captures
	}

	for i, name := range r.SubexpNames(){
		if i == 0{
			continue
		}
		captures[name] = match[i]
	}
	return captures
}

const (
	MongoDBHosts = "www.motrexlab.net:27017"
	AuthDatabase = "seobaksa"
//	AuthUserName = "guest"
//	AuthPassword = "welcome"
//	TestDatabase = "goinggo"
)

// main is the entry point for the application.
func InitMongoDB() *mgo.Session{
	// We need this object to establish a session to our MongoDB.
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{MongoDBHosts},
		Timeout:  60 * time.Second,
		Database: AuthDatabase,
//		Username: AuthUserName,
//		Password: AuthPassword,
	}

	// Create a session which maintains a pool of socket connections
	// to our MongoDB.
	mongoSession, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		log.Critical("CreateSession: %s\n", err)
	}

	// Reads may not be entirely up-to-date, but they will always see the
	// history of changes moving forward, the data read will be consistent
	// across sequential queries in the same session, and modifications made
	// within the session will be observed in following queries (read-your-writes).
	// http://godoc.org/labix.org/v2/mgo#Session.SetMode
	mongoSession.SetMode(mgo.Monotonic, true)

	return mongoSession
}

func InsertMapToMongoDB(threadCnt int, data map[string]string, session *mgo.Session){

//	var waitGroup sync.WaitGroup
//
//	 Perform 5 concurrent queries against the database.
//	waitGroup.Add(threadCnt)

//	for i := 0; i < threadCnt; i++{
		go func(){
			c := session.DB("test").C("gpsDeviceInfo")
			if len(data) != 0 {
				err := c.Insert(data)
				if err != nil {
					log.Critical(err)
				}
			}
//			defer waitGroup.Done()
		}()
//	}
//	go func() {
//		waitGroup.Wait()
//	}()
}