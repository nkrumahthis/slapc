package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
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
	stdin, err := session.StdinPipe()
	if err != nil {
		fmt.Println(err.Error())
	}

	stdout, err := session.StdoutPipe()
	if err != nil {
		fmt.Println(err.Error())
	}

	stderr, err := session.StderrPipe()
	if err != nil {
		fmt.Println(err.Error())
	}

	// Goroutine to handle command output
	go func() {
		io.Copy(os.Stdout, stdout)
	}()
	go func() {
		io.Copy(os.Stderr, stderr)
	}()

	session.Shell()

	// Get list of projects
	_, err = fmt.Fprintln(stdin, "ls -d "+chosenServer.Path+"/*/")
	if err != nil {
		fmt.Println("Failed to send command: ", err)
		return
	}

	// Wait for the commands to finish
	err = session.Wait()
	if err != nil {
		fmt.Println("Failed to wait for session: ", err)
		return
	}

	// session := CreateSSHSession(connection)
}
