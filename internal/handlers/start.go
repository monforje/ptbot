package handlers

import (
	"context"
	"ptbot/internal/service/tgsvc"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	tele "gopkg.in/telebot.v4"
)

func StartHandler(db *mongo.Database) tele.HandlerFunc {
	return func(c tele.Context) error {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if c.Sender() == nil {
			return c.Send("unable to retrieve user information")
		}

		col := db.Collection("users")
		result := tgsvc.Start(ctx, col, c.Sender().ID)
		return c.Send(result)
	}
}
