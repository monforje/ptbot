package handlers

import (
	"ptbot/internal/service/photo"
	"ptbot/internal/service/tags"

	"go.mongodb.org/mongo-driver/v2/mongo"
	tele "gopkg.in/telebot.v4"
)

func UploadHandler(db *mongo.Database) tele.HandlerFunc {
	return photo.Upload(db)
}

func SetNameHandler(db *mongo.Database) tele.HandlerFunc {
	return tags.SetName(db)
}

func AddTagsHandler(db *mongo.Database) tele.HandlerFunc {
	return tags.Generate(db)
}
