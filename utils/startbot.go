package utils

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"telegram/config"

	log "github.com/sirupsen/logrus"

	tb "gopkg.in/tucnak/telebot.v2"
)

var (
	Menu           = &tb.ReplyMarkup{}
	MenuIn         = &tb.ReplyMarkup{}
	BtnNewUser     = Menu.Text("Добавить нового пользователя")
	BtnNewOrigin   = Menu.Text("Подключить новый сервис")
	BtnMyId        = Menu.Text("Мой ID")
	BtnShowOrigins = MenuIn.Data("Подключенные сервисы", "origins")
	BtnAddOrigin   = MenuIn.Data("Добавить новый", "newOrigin")
	AuthClient     = &http.Client{Timeout: 10 * time.Second}
	Bot            = &tb.Bot{}
)

type Recipient struct {
	ID int
}

func (user Recipient) Recipient() string {
	return strconv.Itoa(user.ID)
}

func StartTelegramBot(ctx context.Context) {
	settings := tb.Settings{
		Token: config.Args.TG_BOT_KEY,
		Poller: &tb.LongPoller{
			Timeout: 1 * time.Second,
		},
	}

	Bot, _ = tb.NewBot(settings)

	Bot.Handle("/start", OnStart())
	Bot.Handle(&BtnMyId, ShowMyId())

	Bot.Handle(&BtnNewUser, func(m *tb.Message) {
		log.Info("BtnNewUser clicked")
		userChat, message := GetId(m)
		if message != "" {
			Bot.Send(userChat, message, Menu)
		}
	})

	Bot.Handle(&BtnNewOrigin, NewOrigin())

	go func() {
		Bot.Start()
	}()

	log.Info("Telegram Bot started")

	<-ctx.Done()

	log.Info("Telegram Bot stopped")
	Bot.Stop()
}
