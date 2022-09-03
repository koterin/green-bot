package controller

import (
	"strings"
	"telegram/internal/utils"

	tb "gopkg.in/telebot.v3"
)

func UsersInlineKeyboard(menu *tb.ReplyMarkup) {
	var (
		users []string
		btn   tb.InlineButton
		btns  []tb.InlineButton
	)

	origins, _ := utils.GetOrigins()
	users = strings.Split(origins, "\n")

	inlineKeys := make([][]tb.InlineButton, 0, 0)
	menu.InlineKeyboard = inlineKeys

	for _, user := range users {
		btn = tb.InlineButton{Unique: user, Text: user}
		btns = []tb.InlineButton{btn}
		menu.InlineKeyboard = append(menu.InlineKeyboard, btns)
	}
}
