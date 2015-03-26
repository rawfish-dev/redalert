package pinger

import (
	"encoding/json"
	"io/ioutil"
)

type PingerConfig struct {
	TargetServers []ServerDetails `json:"servers"`
}

type ServerDetails struct {
	Name     string `json:"name"`
	Address  string `json:"address"`
	Interval int    `json:"interval"`
}

func ReadConfigFile() (*PingerConfig, error) {
	file, err := ioutil.ReadFile("pinger/config.json")
	if err != nil {
		return nil, err
	}
	var config PingerConfig
	err = json.Unmarshal(file, &config)
	return &config, err
}
