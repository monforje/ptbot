package tgsvc

import (
	tele "gopkg.in/telebot.v4"
)

func Info(c tele.Context) error {
	msg := `
PicsTagsBot — твой личный архив фото в Telegram

Храни и систематизируй свои изображения быстро и безопасно

Приватность, Контроль доступа, Удобная организация
`

	return c.Send(msg, &tele.SendOptions{
		ParseMode: tele.ModeMarkdown,
	})
}
