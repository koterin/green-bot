package utils

import (
	"context"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	tb "gopkg.in/tucnak/telebot.v2"
)

var (
	menu       = &tb.ReplyMarkup{}
	btnNewUser = menu.Text("Добавить нового пользователя")
	btnMyId    = menu.Text("Мой ID")
)

type Recipient struct {
	ID int
}

func (user Recipient) Recipient() string {
	return strconv.Itoa(user.ID)
}

func StartTelegramBot(ctx context.Context, TG_BOT_KEY string) {
	settings := tb.Settings{
		Token: TG_BOT_KEY,
		Poller: &tb.LongPoller{
			Timeout: 1 * time.Second,
		},
	}

	bot, err := tb.NewBot(settings)
	if err != nil {
		log.Fatal(err)
	}

	menu.Reply(
		menu.Row(btnMyId),
		menu.Row(btnNewUser),
	)

	bot.Handle("/start", func(m *tb.Message) {
		if !m.Private() {
			log.Error("Error: chat is not private")
			return
		}

		log.Info("User started bot: ", m.Sender.Username)

		var userChat Recipient
		userChat.ID = int(m.Chat.ID)

		message := "Сообщи этот ID админу для авторизации: " + userChat.Recipient()
		bot.Send(userChat, message, menu)
	})

	go func() {
		bot.Start()
	}()

	log.Info("Telegram Bot started")

	<-ctx.Done()

	log.Info("Telegram Bot stopped")
	bot.Stop()
}
