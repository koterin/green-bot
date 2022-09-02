package controller

import (
	"context"
	"time"

	"telegram/config"

	log "github.com/sirupsen/logrus"

	tb "gopkg.in/telebot.v3"
)

var Bot = &tb.Bot{}

func StartTelegramBot(ctx context.Context) {
	settings := tb.Settings{
		Token: config.Args.TG_BOT_KEY,
		Poller: &tb.LongPoller{
			Timeout: 1 * time.Second,
		},
	}

	Bot, _ = tb.NewBot(settings)

	Bot.Handle("/start", OnStart())
	/*Bot.Handle(tb.OnText, OnText())
	// Bot.Handle(tb.OnQuery)
	*/
	// Buttons
	Bot.Handle(&BtnMyId, ShowMyId())
	Bot.Handle(&BtnNewUser, NewUser())
	/*Bot.Handle(&BtnNewOrigin, NewOrigin())
	Bot.Handle(&BtnNewPermission, NewPermission())

	// Inline Buttons
	Bot.Handle(&BtnShowOrigins, ShowOrigins())
	Bot.Handle(&BtnAddOrigin, AddOrigin())
	Bot.Handle(&BtnAddUser, AddUser())
	*/
	go func() {
		Bot.Start()
	}()

	log.Info("Telegram Bot started")

	<-ctx.Done()

	log.Info("Telegram Bot stopped")
	Bot.Stop()
}
