package command

import (
	"context"
	"fmt"
	"ptbot/pkg/erro"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func GetByID[T any](ctx context.Context,
	col *mongo.Collection,
	tgid int64,
) (T, error) {
	var result T

	filter := bson.M{"tg_id": tgid}

	err := col.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return result, erro.ErrDocumentNotFound
		}
		return result, fmt.Errorf("failed to find document: %w", err)
	}

	return result, nil
}
