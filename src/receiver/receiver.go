package main

import(
	"fmt"
	"runtime"
	"net"
	"sync"
	"regexp"
	"github.com/spf13/viper"
	"strconv"
)

// 데이터를 슬라이스에 넣을 때는 성능향상을 위해 pool을 이용하는 것을 검토해본다.

// 슬라이스에서 가져온 데이터를 차례대로 파싱하기 위한 함수를 만들어야 한다.
// 1. 데이터 인코딩 2. 정규표현식으로 파싱 3. key, value 형태로 맵에 저장(pool) 이용
// 파싱을 하기 위해서는 정규표현식을 이용하며, 정규표현식은 config에서 가져온다.
// 정규표현식으로 데이터를 Key, value 형태로 분리하고, 데이터 베이스에 저장한다.
// key,value로 쪼개진 데이터를 mongo db에 써줘야하는데, 이때 mysql도 고려한다.

// 고려사항
// CPU 수도 고려
//

const configfilename = "config"
const configpath = "D:\\godev\\nfleetd\\src\\receiver\\."

func main() {
	//use all CPU core
	numCPUs := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPUs)

	//read Configuration file
	viper.SetConfigName(configfilename) // name of config file (without extension)
	viper.AddConfigPath(configpath)      // path to look for the config file in

	err := viper.ReadInConfig()

	result := make(chan string)

	if err != nil {
		fmt.Println("Config not found... ", err)
	} else {
		for _, name := range viper.AllKeys(){

			//config default setting
			//viper.SetDefault(name, )
			//Routine call
			status := viper.GetBool(name+".status")
			protocol := viper.GetString(name+".protocol")
			address := viper.GetString(name+".address")
			port := viper.GetInt(name+".port")
			regex := viper.GetString(name+".regex")
			goRoutine := viper.GetInt(name+".goRoutine")
			buffer := viper.GetInt(name+".buffer")
			if status == true{
				fmt.Println("start to receive data : name="+name+" IP="+address+" port="+strconv.Itoa(port))
				go startReciever(name, protocol, address, port, regex, goRoutine, buffer, result)
			}
		}
		for r := range result{
			fmt.Println(r)
		}
		defer close(result)
	}
}

func startReciever (name string, protocol string, address string, port int, regex string, goRoutine int, buffer int, result chan string) {
	ln, err := net.Listen(protocol, address + ":" + strconv.Itoa(port))
	if err != nil{
		fmt.Println(err)
		return
	}
	defer ln.Close()

	ch := make(chan string)
	done := make(chan struct{})

	var wg sync.WaitGroup
	const numWorkers = 3
	wg.Add(numWorkers)

	for i:=0; i< numWorkers; i++{
		go func(n int) {
			worker(n, done, ch, result)
			wg.Done()
		}(i)
	}

	go func() {
		wg.Wait()
	}()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		defer conn.Close()

		go func(c net.Conn) {
			data := make([]byte, 4096)
			for {
				n, err := c.Read(data)
				if err != nil {
					fmt.Println(err)
					return
				}
				ch <- string(data[:n])
			}
		}(conn)
	}

	defer close(done)
	defer close(ch)

}

func worker(n int, done <-chan struct{}, ch<-chan string, result chan <- string){
	a := make(map[string]string)
	for rawData := range ch{
		select{
		case <-done:
			return
		default :
			a = parseData(rawData);
			result <- fmt.Sprintf("worker: %d, bbb : %+v", n, a)
		}
	}
}

type myRegexp struct {
	*regexp.Regexp
}

func parseData(rawData string) map[string]string{
	s1 := make(map[string]string)
	re1 := myRegexp{regexp.MustCompile("(?P<name>\\w+\\s:\\s\\d)")};
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
