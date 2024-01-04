package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	host := os.Getenv("HOST")
	user := os.Getenv("USERNAME")
	pwd := os.Getenv("PASSWORD")
	fmt.Println(host)
	ConnectToServer(host, user, pwd)
	fmt.Println("successful")
}
