package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"os/exec"
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

type VPN struct {
	Host string `json:"host"`
	Port string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Certificate string `json:"certificate"`
}

type Config struct {
	KnownHosts string   `json:"known_hosts"`
	PrivateKey string   `json:"private_key"`
	Servers    []Server `json:"servers"`
	Vpn		   VPN `json:"vpn"`
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

func checkIfOnVPN(vpn VPN) (bool, error) {
	resp, err := http.Get("https://ifconfig.me")

	if err != nil {
		return false, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	print(string(body) + "\n")

	if string(body) == vpn.Host {
		print("on vpn\n")
		return true, nil
	} else {
		print("off vpn\n")
		return false, nil
	}
}

func getOnVPN() (bool, error) {
	exec.Command("")
	return false, nil
}

func main() {

	config := GetAppConfig()

	print("checking if on vpn...\n")

	onVPN, err := checkIfOnVPN(config.Vpn)
	if err != nil {
		panic(err)
	}

	if !onVPN {
		getOnVPN()
	}

}