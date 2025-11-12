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

		sticker := &tele.Sticker{
			File: tele.File{
				FileID: "CAACAgIAAxkBAAET2MxpFGN1Gh-6VwFM7t6qkNIQAujK8AACFwADTlzSKeZelPOMymwmNgQ",
			},
		}
		c.Send(sticker)

		col := db.Collection("users")
		result := tgsvc.Start(ctx, col, c.Sender().ID)

		if !result.IsRegistered {
			markup := &tele.ReplyMarkup{}
			btnReg := markup.Data("Зарегистрироваться", "reg_button")
			markup.Inline(
				markup.Row(btnReg),
			)
			return c.Send(result.Message, markup)
		}

		return c.Send(result.Message)
	}
}
