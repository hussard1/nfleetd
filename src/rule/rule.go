package rule

import (
	"fmt"
)

type RuleEngine interface {
	Parse([]byte) *Message
}

type Message struct {
	StartByte string
	PacketLen string
	ProtocolNum string
	Datetime string
	SatelliteNum string
	Latitude float64
	Longtitude float64
	Speed string
	Direction string
	RemainByte string
	SerialNum string
	Checksum string
	StopByte string
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