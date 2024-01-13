package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("started main")
	channel := make(chan string)

	go waitC(channel)
	time.Sleep(time.Second * 1)
	fmt.Println("main finished")

	var response string
	response = <- channel

	fmt.Println(response)
	fmt.Println("end")

}

func waitC(channel chan string) {
	fmt.Println("started waitC")
	time.Sleep(time.Second * 3)
	response := "waitC finished"
	channel <- response
}