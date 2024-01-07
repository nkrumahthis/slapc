package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

var appConfig Config

type Server struct {
	Name string `json:"name"`
	Host string `json:"host"`
	Port string `json:"port"`
	User string `json:"user"`
	Pass string `json:"pass"`
	Path string `json:"path"`
}
type Config struct {
	KnownHosts string   `json:"known_hosts"`
	PrivateKey string   `json:"private_key"`
	Servers    []Server `json:"servers"`
}

func init() {
	appConfig = readConfig()
}

func GetAppConfig() Config {
	return appConfig
}

func readConfig() Config {
	configFileBytes, err := os.ReadFile("./config.json")
	if err != nil {
		panic("Error reading config file " + err.Error())
	}

	var config Config

	err = json.Unmarshal(configFileBytes, &config)
	if err != nil {
		panic("error unmarshalling config " + err.Error())
	}

	return config
}

func main() {
	config := GetAppConfig()
}
