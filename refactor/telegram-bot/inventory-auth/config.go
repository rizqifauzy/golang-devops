package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func config(key string) string {
	// Load .env file
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Print("Error loading .env file")
	}
	return os.Getenv(key)
}
