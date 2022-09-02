package controller

import (
	"telegram/internal/entity"
	"telegram/internal/utils"

	log "github.com/sirupsen/logrus"

	tb "gopkg.in/telebot.v3"
)

var (
	Menu             = &tb.ReplyMarkup{ResizeKeyboard: true}
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
		userChat, message := utils.GetId(c.Message())

		if message != "" {
			if err := utils.IsAdmin(userChat.ID); err != nil {
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

func OnText() func(*tb.Message) {
	return func(m *tb.Message) {
		var (
			state string
			msgID int
			msg   string
		)

		if _, userExist := utils.UserStates[m.Chat.ID]; !userExist {
			utils.UserStates[m.Chat.ID] = make(map[string]int)
		}

		for state, msgID = range utils.UserStates[m.Chat.ID] {
			if msgID == m.ID {
				switch state {
				case entity.StateAddOrigin:
					msg = utils.ValidateOrigin(m.Text)

					MenuIn.Inline(
						MenuIn.Row(BtnShowOrigins),
						MenuIn.Row(BtnAddOrigin),
					)
					Bot.Send(m.Chat, msg, MenuIn)
				case entity.StateAddUserEmail:
					msg = "StateAddUserEmail"
					Bot.Send(m.Chat, msg)
				case entity.StateAddUserHost:
					msg = "StateAddUserHost"
					Bot.Send(m.Chat, msg)
				}

				return
			}
		}

		Bot.Send(m.Chat, entity.TextUnknownMsg)
	}
}
