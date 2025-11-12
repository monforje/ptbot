package photo

import (
	"context"
	"log"
	"ptbot/internal/db/command"
	"ptbot/internal/model"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	tele "gopkg.in/telebot.v4"
)

func Upload(db *mongo.Database) func(c tele.Context) error {
	return func(c tele.Context) error {
		if c.Message().Photo == nil {
			return c.Send("Пожалуйста, отправьте фото")
		}

		photo := c.Message().Photo

		photoDoc := model.Photo{
			FileID:       photo.FileID,
			FileUniqueID: photo.UniqueID,
			SendDate:     time.Now(),
			Size:         photo.FileSize,
			Name:         "",
			Tags:         []string{},
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		col := db.Collection("photos")

		// Проверяем, есть ли уже фото с таким FileUniqueID
		existing, err := command.GetPhotoByFileUniqueID(ctx, col, photo.UniqueID)
		if err == nil && existing != nil {
			// Фото уже есть в базе
			return c.Send("Это фото уже сохранено в базе")
		}
		if err != nil {
			// если это не ErrNoDocuments — логируем, но пытаемся вставить
			if err != mongo.ErrNoDocuments {
				log.Printf("warning: failed to check existing photo: %v", err)
			}
		}

		_, err = col.InsertOne(ctx, photoDoc)
		if err != nil {
			log.Printf("failed to save photo metadata: %v", err)
			return c.Send("Ошибка при сохранении фото")
		}

		return c.Send("Фото сохранено!\n\nТеперь добавьте имя: =имя\nИли тэги: +тэг1, тэг2 ответив на это сообщение")
	}
}
