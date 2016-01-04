package main
import (
	"fmt"
	"net"
	"sync"
	"rule"
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

func (server *Server) Bind(worker Worker, device Device, deviceList map[string]int, geofenceList []Geofence, mongoSession *MongoSession, mysqlSession *SQLMapper) {
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
			execute(ch, imeiMap, re, deviceList, geofenceList, mongoSession, mysqlSession)
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

func execute(ch<-chan DataSet, IMEIMap map[net.Conn]string, re rule.RuleEngine, deviceList map[string]int, geofenceList []Geofence, mongoSession *MongoSession, mysqlSession *SQLMapper) {
	for dataSet := range ch {
		msgList := re.Parse(dataSet.dataLength, dataSet.rawdata, dataSet.conn, IMEIMap)
		mongoSession.InsertMessageToMongoDB(msgList)
		InsertDeviceInfo(msgList, deviceList, mysqlSession)
		CheckInOutGeofence(msgList, geofenceList, mysqlSession)
	}
}

func InsertDeviceInfo(msgList []rule.Message, deviceList map[string]int, mysqlSession *SQLMapper){
	if msgList != nil{
		for _, msg := range msgList{
			if _, ok := deviceList[msg.IMEI]; !ok {
				mysqlSession.InsertDeivceInfoToMysql(msg.IMEI);
				deviceList[msg.IMEI] = 0
			}
		}
	}
}