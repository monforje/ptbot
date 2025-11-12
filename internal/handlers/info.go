package handlers

import (
	"ptbot/internal/service/tgsvc"

	tele "gopkg.in/telebot.v4"
)

func InfoHandler() tele.HandlerFunc {
	return func(c tele.Context) error {
		return tgsvc.Info(c)
	}
}
