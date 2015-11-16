package rule

import (
	"fmt"
	"net"
)

type RuleEngine interface {
	Parse(dataLength int, rawdata []byte, conn net.Conn, IMEIMap map[net.Conn]string) []Message
}

type Message struct {
	IMEI string `json:"imei"`
	Time string `json:"time"`
	Location `json:"location"`
	Cell `json:"cell"`
	Status `json:"status"`
	Event `json:"event"`
}

type Location struct{
	Satellitenum int `json:"satellitenum"`
	Latitude float64 `json:"latitude"`
	Longtitude float64 `json:"longtitude"`
	Speed int `json:"speed"`
	Direction string `json:"direction"`
}

type Cell struct{
	Mcc int `json:"mcc"`
	Mnc int `json:"mnc"`
	Lac int `json:"lac"`
	Cellid int `json:"cellid"`
}

type Status struct{
	Gps `json:"gps"`
	Gsm `json:"gsm"`
}

type Gps struct {
	Voltage int `json:"voltage"`
	Power int `json:"power"`
	Acc int `json:"acc"`
	Battery int `json:"battery"`
}

type Gsm struct {
	Strength int `json:"strength"`
}

type Event struct{
	Lowbattery int `json:"lowbattery"`
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