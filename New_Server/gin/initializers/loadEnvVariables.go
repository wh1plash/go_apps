package initializers

import (
	"github.com/joho/godotenv"
	"log"
)

// LoadEnvVariables loads variables from the .env file into the environment
func LoadEnvVariables() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
