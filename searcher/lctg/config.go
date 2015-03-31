package searcher

import (
	"encoding/json"
	"io/ioutil"
)

type LCTGSearchConfig struct {
	Name       string `json:"name"`
	Login      string `json:"login"`
	Password   string `json:"password"`
	SearchPath string `json:"address"`
	Interval   int    `json:"interval"`
}

func ReadConfigFile() (*LCTGSearchConfig, error) {
	file, err := ioutil.ReadFile(PACKAGE_PATH + "/config.json")
	if err != nil {
		return nil, err
	}
	var config LCTGSearchConfig
	err = json.Unmarshal(file, &config)
	return &config, err
}
