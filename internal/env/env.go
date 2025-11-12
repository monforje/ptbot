package env

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	DBURL    string
	BotToken string
}

func New() (*Env, error) {
	e := &Env{}

	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}

	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("BOT_TOKEN is nil")
	}
	e.BotToken = token

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DB_URL is nil")
	}
	e.DBURL = dbURL

	log.Println("env initialized")

	return e, nil
}
