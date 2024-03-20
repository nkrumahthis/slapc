package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
)

func checkIfOnVPN() (bool, error) {
	resp, err := http.Get("https://ifconfig.me")

	if err != nil {
		return false, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	print(string(body) + "\n")

	if string(body) == "192.168.1.1" {
		print("on vpn\n")
		return true, nil
	} else {
		print("off vpn\n")
		return false, nil
	}
}

func getOnVPN() (bool, error) {
	return false, nil
}

func main() {

	print("checking if on vpn...\n")

	onVPN, err := checkIfOnVPN()
	if err != nil {
		panic(err)
	}

	if !onVPN {
		getOnVPN()
	}

}