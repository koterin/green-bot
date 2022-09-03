package controller

import (
	"telegram/config"
	"telegram/internal/entity"
	"telegram/internal/utils"

	log "github.com/sirupsen/logrus"

	tb "gopkg.in/telebot.v3"
)

func OriginsInlineKeyboard(menu *tb.ReplyMarkup) error {
	var (
		btn  tb.InlineButton
		btns []tb.InlineButton
		data entity.ResponseData
	)

	inlineKeys := make([][]tb.InlineButton, 0, 0)
	menu.InlineKeyboard = inlineKeys

	if err := utils.GetStruct(config.Args.ORIGIN_URL, &data); err != nil {
		log.Error(err)

		return err
	}

	for _, host := range data.Origins {
		btn = tb.InlineButton{Unique: host.Origin, Text: host.Origin}
		btns = []tb.InlineButton{btn}
		menu.InlineKeyboard = append(menu.InlineKeyboard, btns)
	}

	return nil
}
