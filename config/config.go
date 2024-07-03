package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var SecretKey string

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	SecretKey = os.Getenv("SECRET_KEY")
	if SecretKey == "" {
		log.Fatal("SECRET_KEY env variable not set")
	}
}