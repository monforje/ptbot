package handlers

import (
	"context"
	"time"

	"ptbot/internal/model"
	"ptbot/internal/service/mdbsvc"

	"go.mongodb.org/mongo-driver/v2/mongo"
	tele "gopkg.in/telebot.v4"
)

func RegHandler(db *mongo.Database) tele.HandlerFunc {
	return func(c tele.Context) error {
		ctx := context.Background()

		now := time.Now()
		user := model.User{
			TgID:      c.Sender().ID,
			Username:  c.Sender().Username,
			FirstName: c.Sender().FirstName,
			LastName:  c.Sender().LastName,
			CreatedAt: now,
			UpdatedAt: now,
		}

		col := db.Collection("users")

		result := mdbsvc.Reg(ctx, col, user)

		return c.Send(result)
	}
}
