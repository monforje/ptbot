package handlers

import (
	"context"
	"fmt"
	"ptbot/internal/db/command"
	"ptbot/internal/model"
	"ptbot/internal/service/mdbsvc"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	tele "gopkg.in/telebot.v4"
)

func RegHandler(db *mongo.Database) tele.HandlerFunc {
	return func(c tele.Context) error {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if c.Sender() == nil {
			return c.Send("Не удалось получить информацию о пользователе")
		}

		// Проверяем, зарегистрирован ли пользователь
		col := db.Collection("users")
		existingUser, err := command.GetByID[model.User](ctx, col, c.Sender().ID)
		if err == nil {
			// Пользователь уже зарегистрирован
			fullName := fmt.Sprintf("%s %s", existingUser.FirstName, existingUser.LastName)
			userInfo := fmt.Sprintf("Вы уже зарегистрированы ✌️\n\n```\nИмя: %s\nНикнейм: @%s\nАйДи: %d\nТелефон: %s\n```",
				fullName, existingUser.Username, existingUser.TgID, existingUser.Phone)
			return c.Send(userInfo, &tele.SendOptions{ParseMode: tele.ModeMarkdown})
		}
		if err != mongo.ErrNoDocuments {
			// Ошибка при проверке, но не "не найдено"
			return c.Send("Ошибка при проверке регистрации")
		}

		if c.Callback() != nil || (c.Message() != nil && c.Message().Contact == nil && c.Text() == "/reg") {
			if c.Callback() != nil {
				c.Respond()
			}

			markup := &tele.ReplyMarkup{
				ResizeKeyboard:  true,
				OneTimeKeyboard: true,
			}

			btnContact := markup.Contact("📱 Поделиться номером телефона")
			markup.Reply(
				markup.Row(btnContact),
			)

			return c.Send("Для регистрации, пожалуйста, поделитесь своим номером телефона:", markup)
		}

		if c.Message() == nil || c.Message().Contact == nil {
			return nil
		}

		if c.Message().Contact.UserID != c.Sender().ID {
			return c.Send("Пожалуйста, отправьте свой собственный номер телефона")
		}

		phone := c.Message().Contact.PhoneNumber
		now := time.Now()
		user := model.User{
			TgID:      c.Sender().ID,
			Username:  c.Sender().Username,
			FirstName: c.Sender().FirstName,
			LastName:  c.Sender().LastName,
			Phone:     phone,
			CreatedAt: now,
			UpdatedAt: now,
		}

		markup := &tele.ReplyMarkup{
			RemoveKeyboard: true,
		}
		processingMsg, _ := c.Bot().Send(c.Recipient(), "Обрабатываю регистрацию...", markup)

		result := mdbsvc.Reg(ctx, col, user, c)

		if processingMsg != nil {
			c.Bot().Delete(processingMsg)
		}

		if result.StickerMsg != nil {
			c.Bot().Delete(result.StickerMsg)
		}

		c.Send(result.Message)

		if result.User != nil {
			fullName := fmt.Sprintf("%s %s", result.User.FirstName, result.User.LastName)
			userInfo := fmt.Sprintf("```\nИмя: %s\nНикнейм: @%s\nАйДи: %d\nТелефон: %s\n```",
				fullName, result.User.Username, result.User.TgID, result.User.Phone)
			c.Send(userInfo, &tele.SendOptions{ParseMode: tele.ModeMarkdown})
		}

		return nil
	}
}
