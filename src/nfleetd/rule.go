package main

import (
	"fmt"

	"goome"
)

type RuleEngine interface {
	Parse()
}

type Message struct {

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
		return new(goome.GoomeStd), nil
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