package controller

import (
	"strings"
	"telegram/internal/entity"
	"telegram/internal/utils"

	log "github.com/sirupsen/logrus"

	tb "gopkg.in/telebot.v3"
)

func NewPermission() tb.HandlerFunc {
	return func(c tb.Context) error {
		log.Info("BtnNewPermission clicked")

		var (
			users []string
			btn   tb.InlineButton
			btns  []tb.InlineButton
		)

		origins, _ := utils.GetOrigins()
		users = strings.Split(origins, "\n")

		for _, user := range users {
			btn = tb.InlineButton{Unique: user, Text: user}
			btns = []tb.InlineButton{btn}
			MenuIn.InlineKeyboard = append(MenuIn.InlineKeyboard, btns)
		}

		log.Debug("btns: ", btns)
		log.Debug("menu: ", MenuIn)

		return c.Send(entity.TextChooseUserMsg, MenuIn)
	}
}
