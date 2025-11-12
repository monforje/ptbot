package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	ID        bson.ObjectID `bson:"_id,omitempty"`
	TgID      int64         `bson:"tg_id"`
	Username  string        `bson:"username"`
	FirstName string        `bson:"first_name"`
	LastName  string        `bson:"last_name"`
	CreatedAt time.Time     `bson:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at"`
}

func (u User) GetTgID() int64 {
	return u.TgID
}
