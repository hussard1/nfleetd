package main
import (
	"fmt"
	"net"
	"sync"
	"rule"
	"gopkg.in/mgo.v2"
	"io"
)

const (
	BUFFER = 4096
)

type Server struct {

}

type DataSet struct {
	dataLength int
	rawdata []byte
	conn net.Conn
}

func (server *Server) Bind(worker Worker, device Device, geofenceList []Geofence, session *mgo.Session, mysqlSession *SQLMapper) {
	address := fmt.Sprintf("%s:%d", device.address, device.port)
	log.Info(fmt.Sprintf("Bind name=%s, address=%s", device.name, address))

	listener, err := net.Listen(device.protocol, address)

	if err != nil {
		log.Error("Cannot bind hostname: ", err)
		return
	}
	defer listener.Close()

	ch := make(chan DataSet, 1)

	imeiMap := make(map[net.Conn]string)

	re, err := rule.CreateRuleEngine(device.rule)
	if err != nil {
		log.Error("Couldn't create the rule engine", err)
		return
	}

	var wg sync.WaitGroup

	wg.Add(worker.thread)

	for i := 0; i < worker.thread; i++ {
		go func(n int) {
			execute(n, ch, imeiMap, device, re, geofenceList, session, mysqlSession)
			wg.Done()
		}(i)
	}

	go func() {
		wg.Wait()
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Error("Cannot accept hostname: ", err)
			continue
		}
		defer func(){
			conn.Close()
			delete(imeiMap, conn)
		}()

		go func(c net.Conn) {
			dataSet := new(DataSet)
			dataSet.conn = c
			dataSet.rawdata = make([]byte, BUFFER)
			for {
				dataSet.dataLength, err = c.Read(dataSet.rawdata)
				if err != nil {
					log.Error("Cannot read from stream: ", err)
				    if err == io.EOF {
						log.Error("detected closed connection", err)
						delete(imeiMap, conn)
					}
					return
				}
				log.Debug("Receive raw data : ", dataSet.rawdata[:dataSet.dataLength])
				ch <- *dataSet
			}
		}(conn)
	}
	defer close(ch)
}

func execute(n int, ch<-chan DataSet, IMEIMap map[net.Conn]string, device Device, re rule.RuleEngine, geofenceList []Geofence, session *mgo.Session, mysqlSession *SQLMapper) {
	for dataSet := range ch {
		msgList := re.Parse(dataSet.dataLength, dataSet.rawdata, dataSet.conn, IMEIMap)
		insertDataToMongoDB(msgList, session)
		CheckInOutGeofence(msgList, geofenceList, mysqlSession)
	}
}

func insertDataToMongoDB(msgList []rule.Message, session *mgo.Session) {
	go func(){
		if msgList != nil {
			for _, msg := range msgList{
				err := session.DB("nfleet").C("gpsdata").Insert(msg)
				if err != nil {
					log.Error("Cannot insert to Mongodb : ", err)
				}
			}
		}
	}()
}