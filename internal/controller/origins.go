package controller

import (
	"telegram/internal/entity"
	"telegram/internal/utils"

	log "github.com/sirupsen/logrus"

	tb "gopkg.in/telebot.v3"
)

func NewOrigin() tb.HandlerFunc {
	return func(c tb.Context) error {
		log.Info("BtnNewOrigin clicked")

		MenuIn.Inline(
			MenuIn.Row(BtnShowOrigins, BtnAddOrigin),
		)

		return c.Send(entity.TextAddOrigin, MenuIn)
	}
}

func ShowOrigins() func(*tb.Callback) {
	return func(c *tb.Callback) {
		log.Info("BtnShowOrigins clicked")

		origins, err := utils.GetOrigins()
		if err != nil {
			log.Error(err)
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
