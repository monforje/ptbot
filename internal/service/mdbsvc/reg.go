package mdbsvc

import (
	"context"
	"errors"
	"ptbot/internal/db/command"
	"ptbot/internal/model"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	tele "gopkg.in/telebot.v4"
)

type RegResult struct {
	Message       string
	StickerMsg    *tele.Message
	User          *model.User
	AlreadyExists bool
}

func Reg(ctx context.Context, col *mongo.Collection, doc model.User, c tele.Context) RegResult {
	var stickerMsg *tele.Message

	err := command.Create(ctx, col, doc)
	if err != nil {
		if errors.Is(err, command.ErrAlreadyExists) {
			existingUser, getErr := command.GetByID[model.User](ctx, col, doc.TgID)
			if getErr == nil {
				return RegResult{
					Message:       "Пользователь уже зарегистрирован",
					StickerMsg:    nil,
					User:          &existingUser,
					AlreadyExists: true,
				}
			}
			return RegResult{
				Message:       "Пользователь уже зарегистрирован",
				StickerMsg:    nil,
				User:          nil,
				AlreadyExists: true,
			}
		}
		return RegResult{
			Message:       "Регистрация не удалась",
			StickerMsg:    nil,
			User:          nil,
			AlreadyExists: false,
		}
	}

	sticker := &tele.Sticker{
		File: tele.File{
			FileID: "CAACAgIAAxkBAAET15dpFELS7HJPrQeVTZJ96hhafk7rIAACcVcAAnBqIEuHSdQDdDCo-TYE",
		},
	}
	stickerMsg, _ = c.Bot().Send(c.Recipient(), sticker)

	time.Sleep(5 * time.Second)

	return RegResult{
		Message:       "Регистрация прошла успешно",
		StickerMsg:    stickerMsg,
		User:          &doc,
		AlreadyExists: false,
	}
}
