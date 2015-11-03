package rule

import (
	"fmt"
)

type RuleEngine interface {
	Parse([]byte) *Message
}

type Message struct {
	PacketLen string
	IMEI string
	CommandType string
	EventCode int
	Latitude float64
	Longtitude float64
	Datetime string
	GPSStatus string
	GSMStatus int
	Speed int
	Direction int
	SatelliteNum int
	HorizontalPositionAccuracy float64
	Altitude int
	Mileage int
	RunTime int
	BaseStationInformation string
	IOPortStatus string
    AnalogInputValue string
	RFID string
	CheckCode string
}

type Rule uint8

const (
	GOOME_STD Rule = iota
	MEITRACK_STD Rule = iota
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
	}

	var r Rule
	return r, fmt.Errorf("not a valid rule engine: %q", rule)
}