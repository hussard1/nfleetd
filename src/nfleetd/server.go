package main
import (
	"fmt"
	"net"
	"sync"
	"rule"
	"gopkg.in/mgo.v2"
	"encoding/json"
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

func (server *Server) Bind(worker Worker, device Device, session *mgo.Session) {
	address := fmt.Sprintf("%s:%d", device.address, device.port)
	log.Info(fmt.Sprintf("Bind name=%s, address=%s", device.name, address))

	listener, err := net.Listen(device.protocol, address)

	if err != nil {
		log.Error("Cannot bind hostname: ", err)
		return
	}
	defer listener.Close()

//	ch := make(chan []byte)
	ch := make(chan DataSet, 1)

	IMEIMap := make(map[net.Conn]string)

	re, err := rule.CreateRuleEngine(device.rule)
	if err != nil {
		log.Error("Couldn't create the rule engine", err)
		return
	}

	var wg sync.WaitGroup

	wg.Add(worker.thread)

	for i := 0; i < worker.thread; i++ {
		go func(n int) {
			execute(n, ch, IMEIMap, device, re, session)
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
		defer conn.Close()

		go func(c net.Conn) {
			dataSet := new(DataSet)
			dataSet.conn = c
			dataSet.rawdata = make([]byte, BUFFER)
			for {
				dataSet.dataLength, err = c.Read(dataSet.rawdata)
				if err != nil {
					log.Error("Cannot read from stream: ", err)
					return
				}
				ch <- *dataSet
			}
		}(conn)
	}
	defer close(ch)
}

func execute(n int, ch<-chan DataSet, IMEIMap map[net.Conn]string, device Device, re rule.RuleEngine, session *mgo.Session) {

	for dataSet := range ch {
		msgList := re.Parse(dataSet.dataLength, dataSet.rawdata, dataSet.conn, IMEIMap)
		b, _ := json.Marshal(msgList)
		log.Debug("Receive data : ", string(b))
//		fmt.Println(msgList)
		InsertMapToMongoDB(msgList, session)
	}
}

func InsertMapToMongoDB(msgList []rule.Message, session *mgo.Session) {
	go func(){
		if msgList != nil {
			for i := 0; i < len(msgList); i++ {
				err := session.DB("test").C("gpsDeviceInfo").Insert(msgList[i])
				fmt.Println(msgList[i])
				if err != nil {
					log.Error("Cannot insert to Mongodb : ", err)
				}
			}
		}
	}()
}