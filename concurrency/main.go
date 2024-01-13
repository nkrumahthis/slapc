package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("started main")

	stream := make(chan int)
	announce := make(chan int)

	go emit(stream, announce) 

	go catch(stream)

	<-announce

	fmt.Println("End of main")
}

func emit(stream chan int, announce chan int) {
	for i := 0; i < 10; i++ {
		stream <- i
		time.Sleep(time.Second * 1)
	}

	close(announce)

}

func catch(stream chan int) {
	for {
		i := <-stream
		fmt.Println(i * 100)
	}
}
