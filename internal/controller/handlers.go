package controller

import (
	"telegram/internal/entity"
	"telegram/internal/utils"

	log "github.com/sirupsen/logrus"

	tb "gopkg.in/telebot.v3"
)

var (
	Menu             = &tb.ReplyMarkup{}
	MenuIn           = &tb.ReplyMarkup{}
	BtnNewUser       = Menu.Text(entity.TextNewUserBtn)
	BtnNewOrigin     = Menu.Text(entity.TextNewOriginBtn)
	BtnNewPermission = Menu.Text(entity.TextNewPermissionBtn)
	BtnMyId          = Menu.Text(entity.TextMyID)
	BtnShowOrigins   = MenuIn.Data(entity.TextCurrentOriginsBtn, "origins")
	BtnAddOrigin     = MenuIn.Data(entity.TextAddOriginBtn, "newOrigin")
	BtnAddUser       = MenuIn.Data(entity.TextAddUserBtn, "newUser")
)

func OnStart() tb.HandlerFunc {
	return func(c tb.Context) error {
		message := utils.GetId(c.Message())

		if message != "" {
			if err := utils.IsAdmin(int(c.Sender().ID)); err != nil {
				log.Info(err)
				Menu.Reply(
					Menu.Row(BtnMyId),
				)
			} else {
				log.Info("Admin user signed in: ", c.Sender().Username)

				Menu.Reply(
					Menu.Row(BtnMyId),
					Menu.Row(BtnNewUser),
					Menu.Row(BtnNewOrigin),
					Menu.Row(BtnNewPermission),
				)
			}
		}

		return c.Send(message, Menu)
	}
}

func OnText() tb.HandlerFunc {
	return func(c tb.Context) error {
		var (
			state string
			msgID int
			msg   string
		)

		if _, userExist := utils.UserStates[c.Chat().ID]; !userExist {
			utils.UserStates[c.Chat().ID] = make(map[string]int)
		}

		for state, msgID = range utils.UserStates[c.Chat().ID] {
			if msgID == c.Message().ID {
				switch state {
				case entity.StateAddOrigin:
					msg = utils.ValidateOrigin(c.Message().Text)

					MenuIn.Inline(
						MenuIn.Row(BtnShowOrigins),
						MenuIn.Row(BtnAddOrigin),
					)

					return c.Send(msg, MenuIn)
				case entity.StateAddUserEmail:
					msg = "StateAddUserEmail"
					return c.Send(msg)
				case entity.StateAddUserHost:
					msg = "StateAddUserHost"
					return c.Send(msg)
				}
			}
		}

		return c.Send(entity.TextUnknownMsg)
	}
}

func OnQuery() tb.HandlerFunc {
	return func(c tb.Context) error {
		c.Send("unknown query")

		return c.Respond()
	}
}
