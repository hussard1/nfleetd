package main

import(
	"fmt"
	"net"
	"time"
	"strconv"
	"sync"
	"math/rand"
)

type Data struct{
	tag string
	buffer []byte
}

func main(){

	senderCnt := 5

	pool := sync.Pool{
		New : func() interface{}{
			data := new(Data)
			data.tag = "new"
			data.buffer = make([]byte, 30)
			return data
		},
	}

	for i :=0 ; i < senderCnt; i++ {
		go func(n int ){
			startSender(n, pool)
		}(i)
	}

	fmt.Scanln()
}

func startSender(n int, pool sync.Pool){

	fmt.Println("sender " + strconv.Itoa(n) + " is start")

	gpsData := pool.Get().(*Data)
	gpsData.buffer = []byte{0x78, 0x78, 0x1F, 0x12, 0x0B, 0x08, 0x1D, 0x11, 0x2E, 0x10, 0xCF, 0x02, 0x7A, 0xC7, 0xEB, 0x0C,
		0x46, 0x58, 0x49, 0x00, 0x14, 0x8F, 0x01, 0xCC, 0x00, 0x03, 0x80, 0x81, 0x0D, 0x0A}

	client, err := net.Dial("tcp", "127.0.0.1:8000")

	if err != nil{
		fmt.Println(err)
		return
	}
	defer client.Close()


	for{
		_, err = client.Write(gpsData.buffer)
		if err != nil{
			fmt.Println(err)
			return
		}

		gpsData.buffer[11] = byte(rand.Intn(256))
		gpsData.buffer[12] = byte(rand.Intn(256))
		gpsData.buffer[13] = byte(rand.Intn(256))
		gpsData.buffer[14] = byte(rand.Intn(256))

		time.Sleep(1 * time.Millisecond)
		gpsData.tag = "used"
		pool.Put(gpsData)
	}

}