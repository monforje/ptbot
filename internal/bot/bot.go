package bot

import (
	"log"
	"ptbot/internal/handlers"
	"ptbot/internal/middleware"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	tele "gopkg.in/telebot.v4"
)

type Bot struct {
	b *tele.Bot
}

func New(token string, db *mongo.Database) (*Bot, error) {
	bot := &Bot{}

	pref := tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		return nil, err
	}

	b.Handle("/reg", handlers.RegHandler(db))
	b.Handle("/start", handlers.StartHandler(db))
	b.Handle(tele.OnContact, handlers.RegHandler(db))
	b.Handle(&tele.Btn{Unique: "reg_button"}, handlers.RegHandler(db))
	b.Handle("/info", handlers.InfoHandler())

	authMiddleware := middleware.RequireRegistration(db)
	b.Handle(tele.OnPhoto, handlers.UploadHandler(db), authMiddleware)
	b.Handle(tele.OnText, func(c tele.Context) error {
		text := c.Message().Text

		if len(text) > 0 && text[0] == '=' {
			return handlers.SetNameHandler(db)(c)
		}

		if len(text) > 0 && text[0] == '+' {
			return handlers.AddTagsHandler(db)(c)
		}

		return nil
	}, authMiddleware)

	bot.b = b

	log.Println("bot initialized")

	return bot, nil
}

func (b *Bot) Run() {
	log.Println("bot started")
	b.b.Start()
}

func (b *Bot) Stop() {
	log.Println("bot stopped")
	b.b.Stop()
}
