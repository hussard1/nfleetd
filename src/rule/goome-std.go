package rule

import (
	"encoding/hex"
//	"strconv"
	"net"
	"fmt"
	"util"

)

type GoomeStd struct {
}

func (re *GoomeStd) Parse(dataLength int, rawdata []byte, conn net.Conn) *Message{
	msg := new(Message)
	fmt.Println(hex.EncodeToString(rawdata[:dataLength]))
	fmt.Println(dataLength)
	fmt.Println(conn)
	responseGoomeLoginData(rawdata, conn)
//	getGPSData(conn)
//	msg = parseGoomeData(raw, msg)
//	msg = calculateData(msg)
	return msg
}

func responseGoomeLoginData(rawdata []byte, conn net.Conn){
	responseLoginData := []byte{0x78, 0x78, 0x05, 0x01, rawdata[12], rawdata[13], byte(util.Crc16(rawdata[14:15])), byte(util.Crc16(rawdata[15:16])), 0x0D, 0x0A}
	_, err := conn.Write(responseLoginData)
	if err != nil{
		fmt.Println("Goome Login data failed to response: ", err)
	}
}

func getGPSData(conn net.Conn){
//	go func(conn net.Conn) {
		data := make([]byte, 100)
		for {
			n, err := conn.Read(data)
			if err != nil {
				fmt.Println("Goome GPS data failed to read: ", err)
				return
			}

		fmt.Println(hex.EncodeToString(data[:n]))
		}
//	}(conn)
}

func parseGoomeData(raw []byte, msg *Message) *Message{
//	data := hex.EncodeToString(raw)
//	msg.StartByte = data[0:4]
//	msg.PacketLen = data[4:6]
//	msg.ProtocolNum = data[6:8]
//	msg.Datetime = data[8:20]
//	msg.SatelliteNum = data[20:22]
//	latitude, _ := strconv.ParseInt(data[22:30], 16, 64)
//	msg.Latitude = float64(latitude)
//	longtitude, _ := strconv.ParseInt(data[30:38], 16, 64)
//	msg.Longtitude = float64(longtitude)
//	msg.Speed = data[38:40]
//	msg.Direction = data[40:44]
//	msg.RemainByte = data[44:48]
//	msg.SerialNum = data[48:52]
//	msg.Checksum = data[52:56]
//	msg.StopByte = data[56:60]
	return msg
}

func calculateData(msg *Message) *Message{
	msg.Latitude = calculateLatitude(msg.Latitude)
	msg.Longtitude = calculateLongtitude(msg.Longtitude)
	return msg
}


func calculateLatitude(lat float64) float64{
	result1 := float64(int(lat/3))/10000
	result2 := int(result1)/60
	return float64(result2)+(result1 - float64(result2*60))/100
}

func calculateLongtitude(long float64) float64{
	result1 := float64(int(long/3))/10000
	result2 := int(result1)/60
	return float64(result2)+(result1 - float64(result2*60))/100
}