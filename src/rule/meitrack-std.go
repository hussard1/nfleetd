package rule

import (
	"strings"
	"strconv"
	"net"
	"bytes"
)

type MeitrackStd struct {
}

func (re *MeitrackStd) Parse(dataLength int, rawdata []byte, conn net.Conn, IMEIMap map[net.Conn]string) []Message{
	msg := new(Message)
	msgList := make([]Message, 0)
	msg = parseMeitrackData(rawdata[:dataLength], msg)
	msgList = append(msgList, *msg)
	return msgList
}

func parseMeitrackData(raw []byte, msg *Message) *Message{
	data := strings.Split(string(raw), ",")
	msg.Devicetype = Meitrack_mvt380
	msg.Messagetype = Location_message
	msg.IMEI = data[1]
//	msg.GPSStatus, _ = strconv.Atoi(data[3])
	msg.Latitude, _ = strconv.ParseFloat(data[4], 64)
	msg.Longitude, _ = strconv.ParseFloat(data[5], 64)
	msg.Time = parseMeitrackDatetimeData(data[6])
//	msg.GPSStatus = data[7]
	msg.Satellitenum, _ = strconv.ParseInt(data[8], 10, 64)
	msg.Strength, _ = strconv.Atoi(data[9])
	msg.Strength = msg.Strength/8
	msg.Speed, _ = strconv.Atoi(data[10])
	msg.Direction, _ = strconv.Atoi(data[11])
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

func parseMeitrackDatetimeData(datetime string) string{

	var buffer bytes.Buffer

	buffer.WriteString("20")
	buffer.WriteString(datetime)

	return buffer.String()
}