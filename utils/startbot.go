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
	Menu         = &tb.ReplyMarkup{}
	BtnNewUser   = Menu.Text("Добавить нового пользователя")
	BtnNewOrigin = Menu.Text("Подключить новый сервис")
	BtnMyId      = Menu.Text("Мой ID")
	AuthClient   = &http.Client{Timeout: 10 * time.Second}
	Bot          = &tb.Bot{}
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
	//if err != nil {
	//	log.Fatal(err)
	//}

	Bot.Handle("/start", OnStart())

	Bot.Handle(&BtnMyId, func(m *tb.Message) {
		log.Info("Button My ID")
		userChat, message := GetId(m)
		if message != "" {
			Bot.Send(userChat, message, Menu)
		}
	})

	Bot.Handle(&BtnNewUser, func(m *tb.Message) {
		log.Info("Button NewUser")
		userChat, message := GetId(m)
		if message != "" {
			Bot.Send(userChat, message, Menu)
		}
	})

	Bot.Handle(&BtnNewOrigin, func(m *tb.Message) {
		log.Info("Button NewOrigin")
		userChat, message := GetId(m)
		if message != "" {
			Bot.Send(userChat, message, Menu)
		}
	})

	go func() {
		Bot.Start()
	}()

	log.Info("Telegram Bot started")

	<-ctx.Done()

	log.Info("Telegram Bot stopped")
	Bot.Stop()
}
