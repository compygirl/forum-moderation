package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Address string `json:"address"`
	DbPath  string `json:"db_path"`
	DbName  string `json:"db_driver"`
}

func CreateConfig() *Config {
	return &Config{}
}

func ReadConfig(configFilepath string, config *Config) error {
	// makes the json file
	configJsonData, err := ioutil.ReadFile(configFilepath)
	if err != nil {
		return err
	}

	// fils the struct config which was defined above
	if err := json.Unmarshal(configJsonData, config); err != nil {
		return err
	}

	return nil
}
