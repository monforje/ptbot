package command

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

var ErrAlreadyExists = errors.New("document already exists")

type TelegramDocument interface {
	GetTgID() int64
}

func Create[T TelegramDocument](ctx context.Context,
	col *mongo.Collection,
	doc T,
) error {
	_, err := GetByID[T](ctx, col, doc.GetTgID())
	if err == nil {
		return ErrAlreadyExists
	}

	if err != mongo.ErrNoDocuments {
		return fmt.Errorf("failed to check document existence: %w", err)
	}

	_, err = col.InsertOne(ctx, doc)
	if err != nil {
		return fmt.Errorf("failed to insert document: %w", err)
	}

	return nil
}
