package controller

import (
	"telegram/internal/entity"

	log "github.com/sirupsen/logrus"

	tb "gopkg.in/telebot.v3"
)

func NewPermission() func(*tb.Message) {
	return func(m *tb.Message) {
		log.Info("BtnNewPermission clicked")

		var users []string
		users = append(users, "user1")
		users = append(users, "user2")

		MenuIn.Inline(
			MenuIn.Row(MenuIn.Data(users[0], users[0])),
		)

		MenuIn.Inline(
			MenuIn.Row(MenuIn.Data(users[1], users[1])),
		)

		Bot.Send(m.Chat, entity.TextChooseUserMsg, MenuIn)
	}
}
