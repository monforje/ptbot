package tgsvc

import (
	"context"
	"fmt"

	"ptbot/internal/db/command"
	"ptbot/internal/model"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

func Start(ctx context.Context, col *mongo.Collection, tgid int64) string {
	user, err := command.GetByID[model.User](ctx, col, tgid)
	if err != nil {
		return "Hi! Please register first with /reg"
	}

	name := user.FirstName
	if name == "" {
		name = user.Username
	}
	if name == "" {
		name = "friend"
	}

	return fmt.Sprintf("Hi, %s", name)
}
