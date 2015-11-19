package rule

import (
	"encoding/hex"
	"net"
	"fmt"
	"util"
	"strconv"
	"bytes"
	"time"
)

const timeformat = "060102150405"
const devicetype = 1

type GoomeStd struct {
}

func (re *GoomeStd) Parse(dataLength int, rawdata []byte, conn net.Conn, IMEIMap map[net.Conn]string) []Message{

	msg := new(Message)

	msgList := make([]Message, 0)

	var startPoint int = -1
	var endPoint int

	for i := 0; i < dataLength; i++{
		if rawdata[i] == 0x78 && rawdata[i+1] == 0x78{
			startPoint = i
		}else if rawdata[i] == 0x0D && rawdata[i+1] == 0x0A{
			endPoint = i+2
			if startPoint != -1 && (endPoint - startPoint) > 14 && (endPoint - startPoint) < 45{
				msg = parseGoomeRawData(rawdata[startPoint:endPoint], msg)
				if value, ok := IMEIMap[conn]; ok {
					msg.IMEI = value
				}else{
					IMEIMap[conn] = msg.IMEI
				}
				msgList = append(msgList, *msg)
			}
		}
	}

	if (dataLength == 15 && rawdata[3] == 0x13) ||
		(dataLength == 18 && rawdata[3] == 0x01) {
		responseGoomeData(rawdata, dataLength, conn)
	}

	return msgList
}

func responseGoomeData(rawdata []byte, dataLength int, conn net.Conn){
	responseLoginData := []byte{0x78, 0x78, 0x05, 0x01, 0x00, 0x01, 0xD9, 0xDC, 0x0D, 0x0A}

	responseLoginData[3] = rawdata[3]
	responseLoginData[4] = rawdata[dataLength-6]
	responseLoginData[5] = rawdata[dataLength-5]
	responseLoginData[6] = byte(util.Crc16(rawdata[dataLength-4:dataLength-3]))
	responseLoginData[7] = byte(util.Crc16(rawdata[dataLength-3:dataLength-2]))

	_, err := conn.Write(responseLoginData)
	if err != nil{
		fmt.Println("Failed to response Goome Login data: ", err)
	}
}

func parseGoomeRawData(rawdata []byte, msg *Message) *Message{

	data := hex.EncodeToString(rawdata)
	datetime := time.Now().Format(timeformat)

	//parse GPS status data
	if len(rawdata) == 15 && rawdata[3] == 0x13{
		msg.Messagetype = 2
		status := strconv.FormatInt(int64(rawdata[4]), 2)
		msg.Acc, _ = strconv.Atoi(status[0:1])
		msg.Power, _ = strconv.Atoi(status[1:2])
		msg.Voltage = int(rawdata[5])
		msg.Strength = int(rawdata[6])
		msg.Time = datetime
	//parse GPS login data
	}else if len(rawdata) == 18 && rawdata[3] == 0x01{
		msg.Messagetype = 1
		msg.IMEI = data[9:24]
		msg.Time = datetime
	//parse GPS location data
	}else if len(rawdata) == 38 && rawdata[3] == 0x12{
		msg.Messagetype = 3
		msg.Time = parseGoomeDatetimeData(rawdata[4:10])
		msg.Location.Satellitenum , _ = strconv.ParseInt(data[21:22], 16, 32)
		Latitude, _ := strconv.ParseInt(data[22:30], 16, 32)
		msg.Latitude = parseGoomeLatitudeData(Latitude)
		Longtitude, _ := strconv.ParseInt(data[30:38], 16, 32)
		msg.Longtitude = parseGoomeLongtitudeData(Longtitude)
		msg.Speed = int(rawdata[20])
		msg.Direction = parseGoomeDirectionData(data[42:46])
	}else if len(rawdata) == 42{

	}
	msg.Devicetype = devicetype
	return msg
}

func parseGoomeDatetimeData(rawdata []byte) string{

	var buffer bytes.Buffer

	for i:=0; i < len(rawdata); i++{
		if i > 0 && rawdata[i] < 10{
			buffer.WriteString("0")
			buffer.WriteString(strconv.Itoa(int(rawdata[i])))

		}else{
			buffer.WriteString(strconv.Itoa(int(rawdata[i])))
		}
	}
	return buffer.String()
}

func parseGoomeDirectionData(data string) int{
	decimaldata, _ := strconv.ParseInt(data, 16, 64)
	binarydata := strconv.FormatInt(decimaldata, 2)
	resultdata, _ := strconv.ParseInt(binarydata[6:], 2, 64)
	return int(resultdata)
}

func parseGoomeLatitudeData(lat int64) float64{
	result1 := lat/3
	result2 := result1/600000
	return float64(result2) + (float64(result1 - result2*600000)/1000000)
}

func parseGoomeLongtitudeData(long int64) float64{
	result1 := long/3
	result2 := result1/600000
	return float64(result2) + (float64(result1 - result2*600000)/1000000)
}
