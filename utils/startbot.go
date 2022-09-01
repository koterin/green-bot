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
	UserStates    = make(map[int64]map[string]int) // map['chatID'] = map'btnAddOrigin' = 'message.ID'
	BackendClient = &http.Client{Timeout: 10 * time.Second}
	Bot           = &tb.Bot{}
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

	Bot.Handle(tb.OnText, OnText())
	Bot.Handle("/start", OnStart())
	Bot.Handle(&BtnMyId, ShowMyId())
	Bot.Handle(&BtnNewUser, NewUser())
	Bot.Handle(&BtnNewOrigin, NewOrigin())

	// Inline:
	Bot.Handle(&BtnShowOrigins, ShowOrigins())
	Bot.Handle(&BtnAddOrigin, AddOrigin())
	Bot.Handle(&BtnAddUser, AddUser())

	go func() {
		Bot.Start()
	}()

	log.Info("Telegram Bot started")

	<-ctx.Done()

	log.Info("Telegram Bot stopped")
	Bot.Stop()
}
