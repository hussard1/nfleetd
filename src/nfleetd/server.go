package main
import (
	"fmt"
	"net"
	"sync"
	"rule"
)

const (
	BUFFER = 4096
)

type Server struct {

}

func (server *Server) Bind(worker Worker, device Device) {
	address := fmt.Sprintf("%s:%d", device.address, device.port)
	log.Info(fmt.Sprintf("Bind name=%s, address=%s", device.name, address))

	listener, err := net.Listen(device.protocol, address)
	if err != nil {
		log.Error("Cannot bind hostname: ", err)
		return
	}
	defer listener.Close()

	ch := make(chan []byte)

	re, err := rule.CreateRuleEngine(device.rule)
	if err != nil {
		log.Error("Couldn't create the rule engine", err)
		return
	}

	var wg sync.WaitGroup
	wg.Add(worker.thread)

	for i := 0; i < worker.thread; i++ {
		go func(n int) {
			execute(n, ch, device, re)
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
			data := make([]byte, BUFFER)
			for {
				n, err := c.Read(data)
				if err != nil {
					log.Error("Cannot read from stream: ", err)
					return
				}
				ch <- data[:n]
			}
		}(conn)
	}
	defer close(ch)
}

func execute(n int, ch<-chan []byte, device Device, re rule.RuleEngine) {

//	data := make(map[string]string)
	for raw := range ch {
		re.Parse(raw)
		fmt.Println(re, raw)
	}
}