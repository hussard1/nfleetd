package rule

import (
	"encoding/hex"
	"strconv"
)

type GoomeStd struct {
}

func (re *GoomeStd) Parse(raw []byte) *Message{
	msg := new(Message)
	msg = parseData(raw, msg)
	msg = calculateData(msg)
	return msg
}

func parseData(raw []byte, msg *Message) *Message{
	data := hex.EncodeToString(raw)
	msg.startByte = data[0:4]
	msg.packetLen = data[4:6]
	msg.protocolNum = data[6:8]
	msg.datetime = data[8:20]
	msg.satelliteNum = data[20:22]
	latitude, _ := strconv.ParseInt(data[22:30], 16, 64)
	msg.latitude = float64(latitude)
	longtitude, _ := strconv.ParseInt(data[30:38], 16, 64)
	msg.longtitude = float64(longtitude)
	msg.speed = data[38:40]
	msg.direction = data[40:44]
	msg.remainByte = data[44:48]
	msg.serialNum = data[48:52]
	msg.checksum = data[52:56]
	msg.stopByte = data[56:60]
	return msg
}

func calculateData(msg *Message) *Message{
	msg.latitude = calculateLatitude(msg.latitude)
	msg.longtitude = calculateLongtitude(msg.longtitude)
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