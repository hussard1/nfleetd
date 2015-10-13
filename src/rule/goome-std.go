package rule

type GoomeStd struct {

}

func (re *GoomeStd) Parse(raw []byte) *Message{
	msg := new(Message)
	msg = parseData(raw, msg)
	msg = calculate(msg)
	return msg
}

func parseData(raw[]byte, msg Message) *Message{
	msg.startByte = string(raw[0:1])
	msg.packetLen = int(raw[2])
	msg.protocolNum = string(raw[3])
	msg.datetime = string(raw[4:9])
	msg.satelliteNum = string(raw[10])
	msg.latitude = float32(raw[11:14])
	msg.longtitude = float32(raw[14:17])
	msg.speed = string(raw[18])
	msg.direction = string(raw[19:20])
	msg.remainByte = string(raw[21:28])
	msg.serialNum = string(raw[29:30])
	msg.checksum = string(raw[31:32])
	msg.stopByte = string(raw[33:34])
	return msg
}

func calculate(msg Message){
	calculateLatitude(msg.latitude)
}

func calculateLatitude(lat float32) float32{
	return lat
}