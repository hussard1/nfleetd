package rule

import (
	"net"
//	"fmt"
)

type AnonymousStd struct {
}

func (re *AnonymousStd) Parse(dataLength int, rawdata []byte, conn net.Conn) []Message{
	msg := new(Message)
	msgList := make([]Message, 0)
	msgList = append(msgList, *msg)

//	if dataLength < 25 {
//		replyAnonymous(conn)
//	}
	//	msg = parseGoomeData(raw, msg)
	//	msg = calculateData(msg)
	return msgList
}

func replyAnonymous(conn net.Conn){
	conn.Write([]byte("LOAD"))
}



func parseAnonymousData(raw []byte, msg *Message) *Message{
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
