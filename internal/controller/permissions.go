package controller

import (
	"telegram/internal/entity"
	"telegram/internal/utils"

	log "github.com/sirupsen/logrus"

	tb "gopkg.in/tucnak/telebot.v2"
)

func NewPermission() func(*tb.Message) {
	return func(m *tb.Message) {
		log.Info("BtnNewPermission clicked")

		utils.AddUserState(m.Chat.ID, entity.StateAddPermission, m.ID+2)

		Bot.Send(m.Chat, entity.TextChooseUserMsg)
	}
}
