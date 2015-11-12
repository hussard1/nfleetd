package rule

import (
	"fmt"
	"net"
)

type RuleEngine interface {
	Parse(dataLength int, rawdata []byte, conn net.Conn, IMEIMap map[net.Conn]string) []Message
}

type Message struct {
	PacketLength int
	IMEI string
	Datetime string
	SatelliteNum int
	Latitude float64
	Longtitude float64
	GPSStatus string
	GSMStatus int
	Speed int
	Direction string
	EventCode string
}

type Rule uint8

const (
	GOOME_STD Rule = iota
	MEITRACK_STD Rule = iota
	ANONYMOUS_STD Rule = iota
)

func CreateRuleEngine(rule string) (RuleEngine, error) {
	r, err := parseRule(rule)
	if err != nil {
		return nil, err
	}

	switch r {
	case GOOME_STD:
		return new(GoomeStd), nil
	case MEITRACK_STD:
		return new(MeitrackStd), nil
	case ANONYMOUS_STD:
		return new(AnonymousStd), nil
	}
	var re RuleEngine
	return re, fmt.Errorf("invalid rule engine type: %q", rule)
}

func parseRule(rule string) (Rule, error) {
	switch rule {
	case "goome/std":
		return GOOME_STD, nil
	case "meitrack/std":
		return MEITRACK_STD, nil
	case "anonymous/std":
		return ANONYMOUS_STD, nil
	}

	var r Rule
	return r, fmt.Errorf("not a valid rule engine: %q", rule)
}