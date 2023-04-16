package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type EnvFields struct {
	PORT              string
	ANILIST_CLIENT_ID string
}

func GetEnv() EnvFields {
	env := EnvFields{}

	err := godotenv.Load(".env")

	if err != nil {
		log.Panicf("Error loading .env file: %q", err)
	}

	env.ANILIST_CLIENT_ID = os.Getenv("ANILIST_CLIENT_ID")

	if envPort := os.Getenv("PORT"); len(envPort) > 0 {
		env.PORT = envPort
	} else {
		envPort = "3001"
	}

	return env
}
