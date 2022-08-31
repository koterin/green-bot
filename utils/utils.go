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
	settings     = tb.Settings{
		Token: config.Args.TG_BOT_KEY,
		Poller: &tb.LongPoller{
			Timeout: 1 * time.Second,
		},
	}
	bot, errBot = tb.NewBot(settings)
)

type Recipient struct {
	ID int
}

func (user Recipient) Recipient() string {
	return strconv.Itoa(user.ID)
}

func StartTelegramBot(ctx context.Context) {
	if errBot != nil {
		log.Fatal(errBot)
	}

	bot.Handle("/start", onStart)

	bot.Handle(&btnMyId, func(m *tb.Message) {
		log.Info("Button My ID")
		userChat, message := GetId(m)
		if message != "" {
			bot.Send(userChat, message, menu)
		}
	})

	bot.Handle(&btnNewUser, func(m *tb.Message) {
		log.Info("Button NewUser")
		userChat, message := GetId(m)
		if message != "" {
			bot.Send(userChat, message, menu)
		}
	})

	bot.Handle(&btnNewOrigin, func(m *tb.Message) {
		log.Info("Button NewOrigin")
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

func onStart(m *tb.Message) error {
	userChat, message := GetId(m)
	if message != "" {
		if err := isAdmin(userChat.ID); err != nil {
			log.Error(err)
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

		return nil
	}

	return fmt.Errorf("Error getting chatId from user %s", m.Sender.Username)
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
	req, err := setAdminRequest(chatId)
	if err != nil {
		log.Error("Error setting request for admin: ", err)

		return err
	}

	resp, err := AuthClient.Do(&req)
	if err != nil {
		log.Error("Error checking admin from Auth Backend: ", err)

		return err
	}

	if resp.StatusCode != http.StatusAccepted {
		log.Error("User is not an admin, ID = ", chatId, " resp code is ", resp.StatusCode)

		return fmt.Errorf("User with ID = %d is not an admin", chatId)
	}

	return nil
}

func setAdminRequest(chatId int) (http.Request, error) {
	body, err := json.Marshal(map[string]string{
		"chat-id": fmt.Sprintf("%d", chatId),
	})
	if err != nil {
		return http.Request{}, fmt.Errorf("Error creating json body: %w", err)
	}

	responseBody := bytes.NewReader(body)

	req, err := http.NewRequest("POST", config.Args.AUTH_URL, responseBody)
	if err != nil {
		return http.Request{}, fmt.Errorf("Error creating request to Backend: %w", err)
	}

	req.Header.Set("X-Green-Origin", "telegram-bot")
	req.Header.Set("Api-Key", config.Args.API_KEY)

	return *req, nil
}
