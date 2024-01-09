package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
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
	if len(os.Args) > 1 {
		args := os.Args[1:]
		switch command := args[0]; command {
		case "help":
			fmt.Println("slapc: see logs and pull code")
		}
	}

	config := GetAppConfig()

	for index, srv := range config.Servers {
		fmt.Printf("[%d] %s\n", index+1, srv.Name)
	}

	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	response := input.Text()
	// response := "1"

	responseIndex, err := strconv.Atoi(response)
	if err != nil {
		panic("invalid response")
	}
	if responseIndex > len(config.Servers) {
		panic("response out of range")
	}

	chosenServer := config.Servers[responseIndex-1]
	connection, err := CreateServerConnection(chosenServer)
	if err != nil {
		panic(err.Error())
	}

	defer connection.Close()

	session, err := connection.NewSession()
	if err != nil {
		fmt.Println("Failed to create session: ", err)
		return
	}
	defer session.Close()

	// set up pipes
	// stdin, err := session.StdinPipe()
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }

	// stdout, err := session.StdoutPipe()
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }

	// stderr, err := session.StderrPipe()
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }

	// Get list of projects
	// directories, err = fmt.Fprintln(stdin, "ls -d "+chosenServer.Path+"/*/")
	res1, err := session.Output("ls -d "+chosenServer.Path+"/*")
	if err != nil {
		fmt.Println("Failed to send command: ", err)
		return
	}

	dirPaths := strings.Split(string(res1), "\n")
	var directories []string
	for _, dirPath := range dirPaths {
		if(dirPath != "" && strings.HasPrefix(dirPath, chosenServer.Path)){
			directories = append(directories, dirPath)
		}
	}

	for idx, dir := range directories {
		directoryNumber := strconv.Itoa(idx + 1)
		fmt.Println(directoryNumber + " " + strings.Split(dir, "/")[4])
	}

}
