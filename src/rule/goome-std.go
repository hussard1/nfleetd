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

type GoomeStd struct {
}

func (re *GoomeStd) Parse(rawdataLength int, rawdata []byte, conn net.Conn, IMEIMap map[net.Conn]string) []Message{

	msg := new(Message)
	msgList := make([]Message, 0)

	for _, data := range bytes.SplitAfter(rawdata[:rawdataLength], []byte{0x0d, 0x0a}){
		if(bytes.HasPrefix(data, []byte{0x78, 0x78})){
			dataLength := len(data)
			/*  check packageLength */
			if(int(data[2]) == dataLength-5){
				/* check CRC */
				if checkSum(data, dataLength){
					msg = parseGoomeRawData(data, msg)
					/* save the imei to map from login message type.
					   because location message type has not imei.
					 */
					if value, ok := IMEIMap[conn]; ok {
						msg.IMEI = value
					}else{
						IMEIMap[conn] = msg.IMEI
					}
					msgList = append(msgList, *msg)
				}
			}
		}
	}

	if (rawdataLength == 15 && rawdata[3] == 0x13) ||
		(rawdataLength == 18 && rawdata[3] == 0x01) {
		responseGoomeData(rawdata, rawdataLength, conn)
	}

	return msgList
}

func checkSum(data []byte, dataLength int) bool{
	i := util.Crc16(data[2:dataLength-4])
	var h, l uint8 = uint8(i>>8), uint8(i&0xff)
	if h == uint8(data[dataLength-4]) && l == uint8(data[dataLength-3]) {
		return true
	}else{
		return false
	}
}


func responseGoomeData(rawdata []byte, dataLength int, conn net.Conn){
	responseLoginData := []byte{0x78, 0x78, 0x05, 0x01, 0x00, 0x01, 0xD9, 0xDC, 0x0D, 0x0A}

	responseLoginData[3] = rawdata[3]
	responseLoginData[4] = rawdata[dataLength-6]
	responseLoginData[5] = rawdata[dataLength-5]

	_, err := conn.Write(responseLoginData)
	if err != nil{
		fmt.Println("Failed to response Goome Login data: ", err)
	}
}

func parseGoomeRawData(rawdata []byte, msg *Message) *Message{

	stringRawdata := hex.EncodeToString(rawdata)
	datetime := time.Now().Format(Timeformat)

	//parse GPS status data
	if len(rawdata) == 15 && rawdata[3] == 0x13{
		msg.Messagetype = Status_message
		status := strconv.FormatInt(int64(rawdata[4]), 2)
		msg.Acc, _ = strconv.Atoi(status[5:6])
		msg.Power, _ = strconv.Atoi(status[4:5])
		msg.Voltage = int(rawdata[5])
		msg.Strength = int(rawdata[6])
		msg.Time = datetime
	//parse GPS login data
	}else if len(rawdata) == 18 && rawdata[3] == 0x01{
		msg.Messagetype = Login_message
		msg.IMEI = stringRawdata[9:24]
		msg.Time = datetime
	//parse GPS location data
	}else if len(rawdata) == 38 && rawdata[3] == 0x12{
		msg.Messagetype = Location_message
		msg.Time = parseGoomeDatetimeData(rawdata[4:10])
		msg.Satellitenum, _ = strconv.ParseInt(stringRawdata[21:22], 16, 32)
		Latitude, _ := strconv.ParseInt(stringRawdata[22:30], 16, 32)
		msg.Latitude = parseGoomeLocationData(Latitude)
		Longtitude, _ := strconv.ParseInt(stringRawdata[30:38], 16, 32)
		msg.Longitude = parseGoomeLocationData(Longtitude)
		msg.Speed = int(rawdata[19])
		msg.Direction = parseGoomeDirectionData(stringRawdata[40:44])
	}else if len(rawdata) == 42{

	}
	msg.Devicetype = Goome_u9
	return msg
}

func parseGoomeDatetimeData(rawdata []byte) string{

	var buffer bytes.Buffer

	for i:=0; i < len(rawdata); i++{
		if i == 0 {
			buffer.WriteString("20")
		}else if rawdata[i] < 10 {
			buffer.WriteString("0")
		}
		buffer.WriteString(strconv.Itoa(int(rawdata[i])))
	}
	return buffer.String()
}

func parseGoomeLocationData(locData int64) float64{
	var result float64
	result = float64(int(float64(locData)/18*10))/1000000
	return result
}


func parseGoomeDirectionData(data string) int{
	decimaldata, _ := strconv.ParseInt(data, 16, 64)
	binarydata := strconv.FormatInt(decimaldata, 2)
	resultdata, _ := strconv.ParseInt(binarydata[6:], 2, 64)
	return int(resultdata)
}