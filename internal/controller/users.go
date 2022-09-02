package controller

import (
	"telegram/internal/entity"
	"telegram/internal/utils"

	log "github.com/sirupsen/logrus"

	tb "gopkg.in/tucnak/telebot.v2"
)

func ShowMyId() func(*tb.Message) {
	return func(m *tb.Message) {
		log.Info("BtnMyId clicked")

		userChat, message := utils.GetId(m)

		if message != "" {
			Bot.Send(userChat, message)
		}
	}
}

func NewUser() func(*tb.Message) {
	return func(m *tb.Message) {
		log.Info("BtnNewUser clicked")

		MenuIn.Inline(
			MenuIn.Row(BtnAddUser),
		)

		Bot.Send(m.Chat, entity.TextAddUser, MenuIn)
	}
}

func AddUser() func(*tb.Callback) {
	return func(c *tb.Callback) {
		log.Info("BtnAddUser clicked")

		utils.AddUserState(c.Message.Chat.ID, entity.StateAddUserEmail, c.Message.ID+2)

		Bot.Send(c.Sender, entity.TextSendEmailMsg)
		Bot.Respond(c, &tb.CallbackResponse{})
	}
}
