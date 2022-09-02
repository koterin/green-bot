package controller

import (
	"telegram/internal/entity"
	"telegram/internal/utils"

	log "github.com/sirupsen/logrus"

	tb "gopkg.in/tucnak/telebot.v2"
)

func NewOrigin() func(*tb.Message) {
	return func(m *tb.Message) {
		log.Info("BtnNewOrigin clicked")

		MenuIn.Inline(
			MenuIn.Row(BtnShowOrigins, BtnAddOrigin),
		)

		Bot.Send(m.Chat, entity.TextAddOrigin, MenuIn)
	}
}

func ShowOrigins() func(*tb.Callback) {
	return func(c *tb.Callback) {
		log.Info("BtnShowOrigins clicked")

		origins, err := utils.GetOrigins()
		if err != nil {
			log.Info(err)
			Bot.Send(c.Sender, entity.TextInternalError)

			return
		}

		MenuIn.Inline(
			MenuIn.Row(BtnShowOrigins, BtnAddOrigin),
		)

		Bot.Send(c.Sender, origins, MenuIn)
		Bot.Respond(c, &tb.CallbackResponse{})
	}
}

func AddOrigin() func(*tb.Callback) {
	return func(c *tb.Callback) {
		log.Info("BtnAddOrigin clicked")

		utils.AddUserState(c.Message.Chat.ID, entity.StateAddOrigin, c.Message.ID+2)

		Bot.Send(c.Sender, entity.TextSendHostMsg)
		Bot.Respond(c, &tb.CallbackResponse{})
	}
}
