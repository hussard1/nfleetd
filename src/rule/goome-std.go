package rule

import (
//	"encoding/hex"
//	"strconv"
)

type GoomeStd struct {
}

func (re *GoomeStd) Parse(raw []byte) *Message{
	msg := new(Message)
//	msg = parseGoomeData(raw, msg)
//	msg = calculateData(msg)
	return msg
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