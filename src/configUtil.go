package main

import (
	"encoding/json"
	"os"
)

var configPrefixes = []string{"~/.config/stocktracker/", "/etc/stocktracker/", "./"}

func ReadConfig(stocks *[]Stock) error {
	// read config file
	for _, prefix := range configPrefixes {
		configFile := prefix + "config.json"
		// check if file exists
		if _, err := os.Stat(configFile); err == nil {
			LoadConfig(configFile, stocks)
			return nil
		}
	}
	return nil
}

func LoadConfig(configFile string, stocks *[]Stock) error {
	// read config file
	jsonFile, err := os.Open(configFile)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	jsonDecoder := json.NewDecoder(jsonFile)
	jsonDecoder.Decode(&stocks)

	return nil
}

func SaveConfig(configFile string, stocks []Stock) error {
	// write config file
	jsonFile, err := os.Create(configFile)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	jsonWriter := json.NewEncoder(jsonFile)
	jsonWriter.Encode(stocks)

	return nil
}
