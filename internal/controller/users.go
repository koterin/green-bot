package controller

import (
	"telegram/internal/entity"
	"telegram/internal/utils"

	log "github.com/sirupsen/logrus"

	tb "gopkg.in/telebot.v3"
)

func ShowMyId() tb.HandlerFunc {
	return func(c tb.Context) error {
		log.Info("BtnMyId clicked")

		message := utils.GetId(c.Message())

		if message != "" {
			return c.Send(message)
		}

		return nil
	}
}

func NewUser() tb.HandlerFunc {
	return func(c tb.Context) error {
		log.Info("BtnNewUser clicked")

		MenuIn.Inline(
			MenuIn.Row(BtnAddUser),
		)

		return c.Send(entity.TextAddUser, MenuIn)
	}
}

func AddUser() tb.HandlerFunc {
	return func(c tb.Context) error {
		log.Info("BtnAddUser clicked")

		utils.AddUserState(c.Chat().ID, entity.StateAddUserEmail, c.Message().ID+2)

		c.Send(entity.TextSendEmailMsg)

		return c.Respond()
	}
}
