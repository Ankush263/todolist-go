package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	dsn := os.Getenv("CONNECTION_STRING")
	if dsn == "" {
		log.Fatal("CONNECTION_STRING is not set")
	}

	// dbConn, err := d
}