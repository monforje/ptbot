package handlers

import (
	"context"
	"ptbot/internal/service/tgsvc"

	"go.mongodb.org/mongo-driver/v2/mongo"
	tele "gopkg.in/telebot.v4"
)

func StartHandler(db *mongo.Database) tele.HandlerFunc {
	return func(c tele.Context) error {
		ctx := context.Background()
		col := db.Collection("users")
		result := tgsvc.Start(ctx, col, c.Sender().ID)
		return c.Send(result)
	}
}
