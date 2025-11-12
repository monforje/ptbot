package mdbsvc

import (
	"context"
	"errors"
	"ptbot/internal/db/command"
	"ptbot/pkg/erro"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

func Reg[T command.TelegramDocument](ctx context.Context, col *mongo.Collection, doc T) string {
	err := command.Create(ctx, col, doc)
	if err != nil {
		if errors.Is(err, erro.ErrDocumentExists) {
			return "user already registered"
		}
		return "registration failed"
	}
	return "registration successful"
}
