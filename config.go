package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/davidmz/mustbe"
)

type TablesCfg map[string]struct {
	Keep    bool              `json:"keep"`
	Clean   bool              `json:"clean"`
	Columns map[string]string `json:"columns"`
}

type Config struct {
	EncryptUUIDs bool      `json:"encryptUUIDs"`
	Tables       TablesCfg `json:"tables"`
}

func loadConfig(fileName string) *Config {
	defer mustbe.Catched(wrapError("cannot read config file %s: %w", fileName))

	data := mustbe.OKVal(ioutil.ReadFile(fileName)).([]byte)
	cfg := &Config{Tables: make(TablesCfg)}
	mustbe.OK(json.Unmarshal(data, cfg))

	return cfg
}
