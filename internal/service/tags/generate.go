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

func Generate(db *mongo.Database) func(c tele.Context) error {
	return func(c tele.Context) error {
		if c.Message().ReplyTo == nil {
			return c.Send("Эта команда должна быть ответом на сообщение с фото")
		}

		if c.Message().ReplyTo.Photo == nil {
			return c.Send("Исходное сообщение должно содержать фото")
		}

		text := c.Message().Text
		if !strings.HasPrefix(text, "+") {
			return c.Send("Формат команды: +тэг1, тэг2, тэг3")
		}

		tagsStr := strings.TrimSpace(strings.TrimPrefix(text, "+"))
		if tagsStr == "" {
			return c.Send("Список тэгов не может быть пустым")
		}

		tagsParts := strings.Split(tagsStr, ",")
		var tags []string
		for _, tag := range tagsParts {
			t := strings.TrimSpace(tag)
			t = strings.TrimPrefix(t, "+")
			t = strings.TrimSpace(t)
			if t != "" {
				tags = append(tags, t)
			}
		}

		if len(tags) == 0 {
			return c.Send("Не удалось распознать тэги")
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

		err = command.AddTagsToPhoto(ctx, col, photoID, tags)
		if err != nil {
			log.Printf("failed to add tags: %v", err)
			return c.Send("Ошибка при добавлении тэгов")
		}

		return c.Send("Тэги добавлены: " + strings.Join(tags, ", "))
	}
}
