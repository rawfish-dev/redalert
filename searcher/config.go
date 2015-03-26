package searcher

import (
	"encoding/json"
	"io/ioutil"
)

type SearchConfig struct {
	TargetServers []ServerDetails `json:"servers"`
}

type ServerDetails struct {
	Name     string `json:"name"`
	Address  string `json:"address"`
	Interval int    `json:"interval"`
}

func ReadConfigFile() (*SearchConfig, error) {
	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		return nil, err
	}
	var config SearchConfig
	err = json.Unmarshal(file, &config)
	return &config, err
}
