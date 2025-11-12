package command

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func UpdatePhoto(ctx context.Context, col *mongo.Collection, id bson.ObjectID, update bson.M) error {
	filter := bson.M{"_id": id}
	updateDoc := bson.M{"$set": update}

	result, err := col.UpdateOne(ctx, filter, updateDoc)
	if err != nil {
		return fmt.Errorf("failed to update document: %w", err)
	}

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func AddTagsToPhoto(ctx context.Context, col *mongo.Collection, id bson.ObjectID, tags []string) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$addToSet": bson.M{"tags": bson.M{"$each": tags}}}

	result, err := col.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to add tags: %w", err)
	}

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func GetPhotoByFileID(ctx context.Context, col *mongo.Collection, fileID string) (bson.M, error) {
	var result bson.M
	filter := bson.M{"file_id": fileID}

	err := col.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, mongo.ErrNoDocuments
		}
		return nil, fmt.Errorf("failed to find photo: %w", err)
	}

	return result, nil
}

func GetPhotoByFileUniqueID(ctx context.Context, col *mongo.Collection, fileUniqueID string) (bson.M, error) {
	var result bson.M
	filter := bson.M{"file_unique_id": fileUniqueID}

	err := col.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, mongo.ErrNoDocuments
		}
		return nil, fmt.Errorf("failed to find photo by unique id: %w", err)
	}

	return result, nil
}
