package tags

import (
	"context"
	"log"
	"ptbot/internal/db/command"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	tele "gopkg.in/telebot.v4"
)

func SetName(db *mongo.Database) func(c tele.Context) error {
	return func(c tele.Context) error {
		if c.Message().ReplyTo == nil {
			return c.Send("Эта команда должна быть ответом на сообщение с фото")
		}

		if c.Message().ReplyTo.Photo == nil {
			return c.Send("Исходное сообщение должно содержать фото")
		}

		text := c.Message().Text
		if !strings.HasPrefix(text, "=") {
			return c.Send("Формат команды: =имя фото")
		}

		photoName := strings.TrimSpace(strings.TrimPrefix(text, "="))
		if photoName == "" {
			return c.Send("Имя фото не может быть пустым")
		}

		fileID := c.Message().ReplyTo.Photo.FileID

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		col := db.Collection("photos")

		photoData, err := command.GetPhotoByFileID(ctx, col, fileID)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return c.Send("Фото не найдено в базе данных")
			}
			log.Printf("failed to find photo: %v", err)
			return c.Send("Ошибка при поиске фото")
		}

		photoID := photoData["_id"].(bson.ObjectID)
		update := bson.M{"name": photoName}

		err = command.UpdatePhoto(ctx, col, photoID, update)
		if err != nil {
			log.Printf("failed to update photo name: %v", err)
			return c.Send("Ошибка при обновлении имени фото")
		}

		return c.Send("Имя фото установлено: " + photoName)
	}
}
