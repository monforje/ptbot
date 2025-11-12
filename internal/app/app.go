package app

import (
	"context"
	"log"
	"ptbot/internal/bot"
	"ptbot/internal/db/mongodb"
	"ptbot/internal/env"
)

type App struct {
	b  *bot.Bot
	db *mongodb.Mongo
}

func New() (*App, error) {
	a := &App{}

	env, err := env.New()
	if err != nil {
		return nil, err
	}

	db, err := mongodb.New(env.DBURL)
	if err != nil {
		return nil, err
	}
	a.db = db

	bot, err := bot.New(env.BotToken, db.DB)
	if err != nil {
		return nil, err
	}
	a.b = bot

	log.Println("app initialized")

	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	log.Println("app started")
	go a.b.Run()

	<-ctx.Done()
	return nil
}

func (a *App) Stop() error {
	log.Println("app stopped")
	a.b.Stop()
	a.db.Stop()
	return nil
}
