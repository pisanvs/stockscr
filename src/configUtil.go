package main

import (
	"encoding/json"
	"os"
)

var configPrefixes []string = []string{"~/.config/stocktracker/", "/etc/stocktracker/", "../"}

func ReadConfig(config Config) error {
	// read config file
	for _, prefix := range configPrefixes {
		configFile := prefix + "config.json"
		if _, err := os.Stat(configFile); err != nil {
			return err
		}
		// config file exists
		LoadConfig(configFile, config)
	}

	return nil
}

func LoadConfig(configFile string, config Config) error {
	// read config file
	jsonFile, err := os.Open(configFile)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	jsonDecoder := json.NewDecoder(jsonFile)
	jsonDecoder.Decode(&config)

	return nil
}

func SaveConfig(configFile string, config Config) error {
	// write config file
	jsonFile, err := os.Create(configFile)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	jsonWriter := json.NewEncoder(jsonFile)
	jsonWriter.Encode(config)

	return nil
}
