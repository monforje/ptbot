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

		existing, err := command.GetPhotoByFileUniqueID(ctx, col, photo.UniqueID)
		if err == nil && existing != nil {
			return c.Send("Вы уже загружали это фото")
		}
		if err != nil {
			if err != mongo.ErrNoDocuments {
				log.Printf("warning: failed to check existing photo: %v", err)
			}
		}

		_, err = col.InsertOne(ctx, photoDoc)
		if err != nil {
			log.Printf("failed to save photo metadata: %v", err)
			return c.Send("Ошибка при сохранении фото")
		}

		sticker := &tele.Sticker{
			File: tele.File{
				FileID: "CAACAgIAAxkBAAET2MppFGMLH2LvK3ROmo1kqa84X59_IAACGAADTlzSKT5q3R61ronZNgQ",
			},
		}
		c.Send(sticker)

		return c.Send("Фото успешно загружено", &tele.SendOptions{ParseMode: tele.ModeMarkdown})
	}
}
