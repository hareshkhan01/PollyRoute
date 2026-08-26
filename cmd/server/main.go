package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Hello Pollution.")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed to .env!")
	}
}
