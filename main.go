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

	// outputChannel := make(chan []byte)
	wr := make(chan []byte)
	finish := make(chan bool)

	go watchStdin(stdin, wr)

	go watchStdout(stdout, wr)

	go watchStderr(stderr, wr)

	session.Shell()

	commands := []string{
		"ls -d " + chosenServer.Path + "/*/",
		"ls -d " + chosenServer.Path + "/*/",
	}

	go executeCommands(commands, wr, finish)

	<- finish

	fmt.Println("End of main")
}

func executeCommands(commands []string, wr chan []byte, finish chan bool){
	defer close(wr)

	for i, command := range commands {
		formattedCommand := "echo @[" + strconv.Itoa(i) + "] && " + command + " && echo ^[" + strconv.Itoa(i) + "] \n"
		wr <- []byte(formattedCommand)
		fmt.Println("command sent")
	}
}

func watchStdin(stdin io.WriteCloser, wr chan []byte) {
	for {
		select {
		case d := <-wr:
			_, err := stdin.Write(d)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}
}

func watchStdout(stdout io.Reader, wr chan []byte) {
	scanner := bufio.NewScanner(stdout)
	for {
		if tkn := scanner.Scan(); tkn {
			rcv := scanner.Bytes()
			raw := make([]byte, len(rcv))
			copy(raw, rcv)
			fmt.Println(string(raw))
		} else {
			if scanner.Err() != nil {
				fmt.Println(scanner.Err())
			} else {
				fmt.Println("io.EOF")
			}
			return
		}
	}
}

func watchStderr(stderr io.Reader, wr chan []byte) {
	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}

func watchStdoutForResponse(stdout io.Reader, outputChannel chan []byte) {
	scanner := bufio.NewScanner(stdout)
	for {
		if tkn := scanner.Scan(); tkn {
			fmt.Println("getting response")
			rcv := scanner.Bytes()
			raw := make([]byte, len(rcv))
			copy(raw, rcv)
			fmt.Println("raw" + string(raw))
			outputChannel <- raw
		} else {
			if scanner.Err() != nil {
				fmt.Println(scanner.Err())
			} else {
				fmt.Println("io.EOF")
			}
			return
		}
	}
}
