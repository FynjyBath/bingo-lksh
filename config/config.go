package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	ApiKey    string `json:"apiKey"`
	ApiSecret string `json:"apiSecret"`
	ContestID int    `json:"contestID"`
}

func LoadConfig(filename string) Config {
	var config Config
	configFile, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	if err != nil {
		log.Fatal(err)
	}
	err = configFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	return config
}
