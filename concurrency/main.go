package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("started c")
	var channel chan string
	go waitC(channel)
	time.Sleep(time.Second * 3)
	var response string
	response = <- channel
	fmt.Println(response)
	fmt.Println("end")

}

func waitC(channel chan string) {
	time.Sleep(time.Second * 3)
	fmt.Println("C")
	response := "ok done"
	channel <- response
}