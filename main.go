package main

import (
	"github.com/deliowales/go-technical-test/finnhub-app/cmd"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err)
		return
	}
	cmd.Execute()
}
