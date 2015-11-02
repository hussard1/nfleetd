package rule

import (
	"strings"
	"strconv"
)

type MeitrackStd struct {
}

func (re *MeitrackStd) Parse(raw []byte) *Message{
	msg := new(Message)
	msg = parseMeitrackData(raw, msg)
	return msg
}

func parseMeitrackData(raw []byte, msg *Message) *Message{
	data := strings.Split(string(raw), ",")
	msg.PacketLen = data[0]
	msg.IMEI = data[1]
	msg.CommandType = data[2]
	msg.EventCode, _ = strconv.Atoi(data[3])
	msg.Latitude, _ = strconv.ParseFloat(data[4], 64)
	msg.Longtitude, _ = strconv.ParseFloat(data[5], 64)
	msg.Datetime = data[6]
	msg.GPSStatus = data[7]
	msg.SatelliteNum, _ = strconv.Atoi(data[8])
	msg.GSMStatus, _ = strconv.Atoi(data[9])
	msg.Speed, _ = strconv.Atoi(data[10])
	msg.Direction, _ = strconv.Atoi(data[11])
	msg.HorizontalPositionAccuracy, _ = strconv.ParseFloat(data[12], 64)
	msg.Altitude, _ = strconv.Atoi(data[13])
	msg.Mileage = data[14]
	msg.RunTime = data[15]
	msg.BaseStationInformation = data[16]
	msg.IOPortStatus = data[17]
	msg.AnalogInputValue = data[18]
	msg.RFID = data[19]
	msg.CheckCode = data[20]

	return msg
}