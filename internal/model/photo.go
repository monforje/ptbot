package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Photo struct {
	ID           bson.ObjectID `bson:"_id,omitempty"`
	FileID       string        `bson:"file_id"`
	FileUniqueID string        `bson:"file_unique_id"`
	SendDate     time.Time     `bson:"send_date"`
	Size         int64         `bson:"size"`
	Name         string        `bson:"name"`
}
