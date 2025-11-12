package tgsvc

import (
	"context"
	"fmt"

	"ptbot/internal/db/command"
	"ptbot/internal/model"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type StartResponse struct {
	Message      string
	IsRegistered bool
}

func Start(ctx context.Context, col *mongo.Collection, tgid int64) StartResponse {
	user, err := command.GetByID[model.User](ctx, col, tgid)

	if err != nil {
		return StartResponse{
			Message:      "Привет, друг\n\nЧтобы начать пользоваться ботом, пожалуйста, зарегистрируйся\n\nИспользуй команду /reg или нажми на кнопку ниже\n\nУзнать больше о боте: /info",
			IsRegistered: false,
		}
	}

	name := user.FirstName
	if name == "" {
		name = user.Username
	}
	if name == "" {
		name = "друг"
	}

	return StartResponse{
		Message:      fmt.Sprintf("Привет, %s\n\nТы уже зарегистрирован!\n\nУзнать больше о боте: /info", name),
		IsRegistered: true,
	}
}
