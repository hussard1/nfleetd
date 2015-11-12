package rule

import (
	"strings"
	"strconv"
	"net"
	"fmt"
)

type MeitrackStd struct {
}

func (re *MeitrackStd) Parse(dataLength int, rawdata []byte, conn net.Conn) []Message{
	msg := new(Message)
	msgList := make([]Message, 0)
	msg = parseMeitrackData(rawdata[:dataLength], msg)
	msgList = append(msgList, *msg)
	fmt.Println(string(rawdata[:dataLength]))
	return msgList
}

func parseMeitrackData(raw []byte, msg *Message) *Message{
	data := strings.Split(string(raw), ",")
	msg.PacketLength, _ = strconv.Atoi(data[0][2:len(data[0])])
	msg.IMEI = data[1]
//	msg.GPSStatus, _ = strconv.Atoi(data[3])
	msg.Latitude, _ = strconv.ParseFloat(data[4], 64)
	msg.Longtitude, _ = strconv.ParseFloat(data[5], 64)
	msg.Datetime = data[6]
	msg.GPSStatus = data[7]
	msg.SatelliteNum, _ = strconv.Atoi(data[8])
	msg.GSMStatus, _ = strconv.Atoi(data[9])
	msg.Speed, _ = strconv.Atoi(data[10])
	msg.Direction = data[11]
//	msg.HorizontalPositionAccuracy, _ = strconv.ParseFloat(data[12], 64)
//	msg.Altitude, _ = strconv.Atoi(data[13])
//	msg.Mileage, _ = strconv.Atoi(data[14])
//	msg.RunTime, _ = strconv.Atoi(data[15])
//	msg.BaseStationInformation = data[16]
//	msg.IOPortStatus = data[17]
//	msg.AnalogInputValue = data[18]
//	msg.RFID = data[19]
//	msg.CheckCode = data[20][0:3]

	return msg
}