package command

import (
	"context"
	"fmt"
	"ptbot/pkg/erro"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type TelegramDocument interface {
	GetTgID() int64
}

func Create[T TelegramDocument](ctx context.Context,
	col *mongo.Collection,
	doc T,
) error {
	_, err := GetByID[T](ctx, col, doc.GetTgID())
	if err == nil {
		return erro.ErrDocumentExists
	}

	if err != erro.ErrDocumentNotFound {
		return fmt.Errorf("failed to check document existence: %w", err)
	}

	_, err = col.InsertOne(ctx, doc)
	if err != nil {
		return fmt.Errorf("failed to insert document: %w", err)
	}

	return nil
}
