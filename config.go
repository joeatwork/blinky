package main

import (
	"encoding/json"
	"io/ioutil"
	"time"
)

type Query struct {
	Name   string `json:"Name"`
	Token  string `json:"Token"`
	Secret string `json:"Secret"`
	Event  string `json:"Event"`
	Where  string `json:"Where,omitempty"`
}

type Device struct {
	Device        string `json:"Device"`
	DeviceAddress string `json:"DeviceAddress"`
}

type Configuration struct {
	ServicePort        string        `json:"ServicePort"`
	PollingRateSeconds time.Duration `json:"PollingRateSeconds"`
	ColorLog           string        `json:"ColorLog"`
	Device             Device        `json:"Device"`
	RedQuery           Query         `json:RedQuery`
	GreenQuery         Query         `json:RedQuery`
	BlueQuery          Query         `json:RedQuery`
}

func readConfig(path string) (config Configuration, err error) {
	var configString []byte
	configString, err = ioutil.ReadFile(path)
	if err != nil {
		return
	}
	err = json.Unmarshal(configString, &config)
	return
}
