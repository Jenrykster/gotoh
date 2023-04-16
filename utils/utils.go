package utils

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

const TOKEN_FILE_PATH = "./.token"

type EnvFields struct {
	PORT              string
	ANILIST_CLIENT_ID string
	ANILIST_SECRET    string
	USER_TOKEN        string
}

func GetEnv() EnvFields {
	env := EnvFields{}

	err := godotenv.Load(".env")

	if err != nil {
		log.Panicf("Error loading .env file: %q", err)
	}

	env.ANILIST_CLIENT_ID = os.Getenv("ANILIST_CLIENT_ID")
	env.ANILIST_SECRET = os.Getenv("ANILIST_SECRET")
	token, err := os.ReadFile(TOKEN_FILE_PATH)

	if err != nil {
		if !strings.Contains(err.Error(), "no such file or directory") {
			env.USER_TOKEN = ""
		} else {
			log.Fatalf("Couldn't read token %q", err)
		}
	} else {
		env.USER_TOKEN = string(token)
	}

	if envPort := os.Getenv("PORT"); len(envPort) > 0 {
		env.PORT = envPort
	} else {
		envPort = "3001"
	}

	return env
}
