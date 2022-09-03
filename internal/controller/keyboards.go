package controller

import (
	"telegram/internal/utils"

	log "github.com/sirupsen/logrus"

	tb "gopkg.in/telebot.v3"
)

func OriginsInlineKeyboard(menu *tb.ReplyMarkup) error {
	var (
		btn  tb.InlineButton
		btns []tb.InlineButton
	)

	inlineKeys := make([][]tb.InlineButton, 0, 0)
	menu.InlineKeyboard = inlineKeys

	origins, err := utils.GetOrigins()
	if err != nil {
		log.Error(err)

		return err
	}

	for _, host := range origins {
		btn = tb.InlineButton{Unique: host.Origin, Text: host.Origin}
		btns = []tb.InlineButton{btn}
		menu.InlineKeyboard = append(menu.InlineKeyboard, btns)
	}

	return nil
}
