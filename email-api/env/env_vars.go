package env

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadVars() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

}
