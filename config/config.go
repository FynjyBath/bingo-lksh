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
	Link      string `json:"link"`
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
	log.Println(config.ApiKey)
	log.Println(config.ApiSecret)
	log.Println(config.ContestID)
	log.Println(config.Link)
	return config
}
