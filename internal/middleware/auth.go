package middleware

import (
	"context"
	"ptbot/internal/db/command"
	"ptbot/internal/model"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	tele "gopkg.in/telebot.v4"
)

func RequireRegistration(db *mongo.Database) tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			if c.Sender() == nil {
				return c.Send("Не удалось получить информацию о пользователе")
			}

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			col := db.Collection("users")

			_, err := command.GetByID[model.User](ctx, col, c.Sender().ID)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					return c.Send("Для использования этой функции необходимо зарегистрироваться.\nИспользуйте команду /reg")
				}
				return c.Send("Ошибка при проверке регистрации")
			}

			return next(c)
		}
	}
}
