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
	BUFFER = 300
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

	re, err := rule.CreateRuleEngine(device.rule)
	if err != nil {
		log.Error("Couldn't create the rule engine", err)
		return
	}

	var wg sync.WaitGroup

	wg.Add(worker.thread)

	for i := 0; i < worker.thread; i++ {
		go func(n int) {
			execute(n, ch, device, re, session)
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

func execute(n int, ch<-chan DataSet, device Device, re rule.RuleEngine, session *mgo.Session) {

	for dataSet := range ch {
		msg := re.Parse(dataSet.dataLength, dataSet.rawdata, dataSet.conn)
		b, _ := json.Marshal(msg)
		fmt.Println("Raw data : ", string(dataSet.rawdata[:dataSet.dataLength]))
		log.Debug("Receive data : ", string(b))
		InsertMapToMongoDB(msg, session)
	}
}

func InsertMapToMongoDB(msg *rule.Message, session *mgo.Session) {

	go func() {
		if msg != nil {
			err := session.DB("test").C("gpsDeviceInfo").Insert(msg)
			if err != nil {
				log.Error("Cannot insert to Mongodb : ", err)
			}
		}
	}()
}