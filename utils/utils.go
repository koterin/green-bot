package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"telegram/config"

	log "github.com/sirupsen/logrus"

	tb "gopkg.in/tucnak/telebot.v2"
)

var (
	menu         = &tb.ReplyMarkup{}
	btnNewUser   = menu.Text("Добавить нового пользователя")
	btnNewOrigin = menu.Text("Подключить новый сервис")
	btnMyId      = menu.Text("Мой ID")
	AuthClient   = &http.Client{Timeout: 10 * time.Second}
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

	bot, err := tb.NewBot(settings)
	if err != nil {
		log.Fatal(err)
	}

	bot.Handle("/start", func(m *tb.Message) {
		userChat, message := GetId(m)
		if message != "" {
			if err := isAdmin(userChat.ID); err != nil {
				menu.Reply(
					menu.Row(btnMyId),
				)
			} else {
				menu.Reply(
					menu.Row(btnMyId),
					menu.Row(btnNewUser),
					menu.Row(btnNewOrigin),
				)
			}

			bot.Send(userChat, message, menu)
		}

	})

	bot.Handle(&btnMyId, func(m *tb.Message) {
		userChat, message := GetId(m)
		if message != "" {
			bot.Send(userChat, message, menu)
		}
	})

	bot.Handle(&btnNewUser, func(m *tb.Message) {
		userChat, message := GetId(m)
		if message != "" {
			bot.Send(userChat, message, menu)
		}
	})

	bot.Handle(&btnNewOrigin, func(m *tb.Message) {
		userChat, message := GetId(m)
		if message != "" {
			bot.Send(userChat, message, menu)
		}
	})

	go func() {
		bot.Start()
	}()

	log.Info("Telegram Bot started")

	<-ctx.Done()

	log.Info("Telegram Bot stopped")
	bot.Stop()
}

func GetId(m *tb.Message) (Recipient, string) {
	var userChat Recipient

	if !m.Private() {
		log.Error("Error: chat is not private")
		return Recipient{}, ""
	}

	log.Info("User started bot: ", m.Sender.Username)

	userChat.ID = int(m.Chat.ID)

	message := "Сообщи этот ID админу для авторизации: " + userChat.Recipient()

	return userChat, message
}

func isAdmin(chatId int) error {
	body, err := json.Marshal(map[string]string{
		"userChatId": fmt.Sprintf("%d", chatId),
	})
	if err != nil {
		return fmt.Errorf("Error creating json body: %w", err)
	}

	responseBody := bytes.NewReader(body)

	resp, err := AuthClient.Post(config.Args.AUTH_URL, "application/json", responseBody)
	if err != nil {
		log.Error("Error checking admin from Auth Backend: ", err)

		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("User with ID = %d is not an admin", chatId)
	}

	return nil
}
