package rule

import (
	"fmt"
)

type RuleEngine interface {
	Parse([]byte) *Message
}

type Message struct {
	startByte string
	packetLen string
	protocolNum string
	datetime string
	satelliteNum string
	latitude float64
	longtitude float64
	speed string
	direction string
	remainByte string
	serialNum string
	checksum string
	stopByte string
}

type Rule uint8

const (
	GOOME_STD Rule = iota
)

func CreateRuleEngine(rule string) (RuleEngine, error) {
	r, err := parseRule(rule)
	if err != nil {
		return nil, err
	}

	switch r {
	case GOOME_STD:
		return new(GoomeStd), nil
	}

	var re RuleEngine
	return re, fmt.Errorf("invalid rule engine type: %q", rule)
}

func parseRule(rule string) (Rule, error) {
	switch rule {
	case "goome/std":
		return GOOME_STD, nil
	}

	var r Rule
	return r, fmt.Errorf("not a valid rule engine: %q", rule)
}