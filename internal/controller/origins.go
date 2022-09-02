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

func ShowOrigins() tb.HandlerFunc {
	return func(c tb.Context) error {
		log.Info("BtnShowOrigins clicked")

		origins, err := utils.GetOrigins()
		if err != nil {
			log.Error(err)

			c.Send(entity.TextInternalError)

			return c.Respond()
		}

		MenuIn.Inline(
			MenuIn.Row(BtnShowOrigins, BtnAddOrigin),
		)

		c.Send(origins, MenuIn)

		return c.Respond()
	}
}

func AddOrigin() tb.HandlerFunc {
	return func(c tb.Context) error {
		log.Info("BtnAddOrigin clicked")

		utils.AddUserState(c.Chat().ID, entity.StateAddOrigin, c.Message().ID+2)

		c.Send(entity.TextSendHostMsg)

		return c.Respond()
	}
}
