package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/davidmz/mustbe"
)

type Rules map[string]struct {
	Action  string            `json:"action"`
	Columns map[string]string `json:"columns"`
}

func loadRules(fileName string) Rules {
	defer mustbe.Catched(wrapError("cannot read file %s: %w", fileName))

	data := mustbe.OKVal(ioutil.ReadFile(fileName)).([]byte)
	rules := make(Rules)
	mustbe.OK(json.Unmarshal(data, &rules))

	return rules
}
