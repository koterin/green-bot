package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"telegram/config"

	log "github.com/sirupsen/logrus"

	tb "gopkg.in/tucnak/telebot.v2"
)

func OnStart() func(*tb.Message) {
	return func(m *tb.Message) {
		userChat, message := GetId(m)
		log.Debug(userChat.ID, " ", message)

		if message != "" {
			if err := isAdmin(userChat.ID); err != nil {
				log.Info(err)
				Menu.Reply(
					Menu.Row(BtnMyId),
				)
			} else {
				log.Info("Admin user signed in: ", m.Sender.Username)

				Menu.Reply(
					Menu.Row(BtnMyId),
					Menu.Row(BtnNewUser),
					Menu.Row(BtnNewOrigin),
				)
			}

			Bot.Send(userChat, message, Menu)
		}
	}
}

func NewOrigin() func(*tb.Message) {
	return func(m *tb.Message) {
		log.Info("BtnNewOrigin clicked")

		MenuIn.Inline(
			MenuIn.Row(BtnShowOrigins, BtnAddOrigin),
		)

		Bot.Send(m.Chat, "hi", MenuIn)
	}
}

func ShowMyId() func(*tb.Message) {
	return func(m *tb.Message) {
		log.Info("BtnMyId clicked")

		Bot.Send(m.Chat, "hi", MenuIn)
	}
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
