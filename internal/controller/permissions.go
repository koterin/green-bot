package controller

import (
	"strconv"
	"telegram/internal/entity"
	"telegram/internal/utils"

	log "github.com/sirupsen/logrus"

	tb "gopkg.in/telebot.v3"
)

func NewPermission() tb.HandlerFunc {
	return func(c tb.Context) error {
		log.Info("BtnNewPermission clicked")

		if err := UsersInlineKeyboard(MenuIn); err != nil {
			c.Send(entity.TextInternalError)
		}

		utils.AddUserState(c.Chat().ID, entity.StateAddUserHost, c.Message().ID+1)

		return c.Send(entity.TextChooseUserMsg+strconv.Itoa(c.Message().ID), MenuIn)
	}
}
