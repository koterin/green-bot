package controller

import (
	"telegram/internal/entity"
	"telegram/internal/utils"

	log "github.com/sirupsen/logrus"

	tb "gopkg.in/tucnak/telebot.v2"
)

var (
	Menu           = &tb.ReplyMarkup{ResizeReplyKeyboard: true}
	MenuIn         = &tb.ReplyMarkup{}
	BtnNewUser     = Menu.Text(entity.TextNewUserBtn)
	BtnNewOrigin   = Menu.Text(entity.TextNewOriginBtn)
	BtnMyId        = Menu.Text(entity.TextMyID)
	BtnShowOrigins = MenuIn.Data(entity.TextCurrentOriginsBtn, "origins")
	BtnAddOrigin   = MenuIn.Data(entity.TextAddOriginBtn, "newOrigin")
	BtnAddUser     = MenuIn.Data(entity.TextAddUserBtn, "newUser")
)

func OnStart() func(*tb.Message) {
	return func(m *tb.Message) {
		userChat, message := utils.GetId(m)

		if message != "" {
			if err := utils.IsAdmin(userChat.ID); err != nil {
				log.Info(err)
				Menu.Reply(
					Menu.Row(BtnMyId),
				)
			} else {
				log.Info("Admin user signed in: ", m.Sender.Username)

				Menu.Reply(
					Menu.Row(BtnMyId),
					Menu.Row(BtnNewUser),
					Menu.Row(BtnNewOrigin),
				)
			}

			Bot.Send(userChat, message, Menu)
		}
	}
}

func NewOrigin() func(*tb.Message) {
	return func(m *tb.Message) {
		log.Info("BtnNewOrigin clicked")

		MenuIn.Inline(
			MenuIn.Row(BtnShowOrigins, BtnAddOrigin),
		)

		Bot.Send(m.Chat, entity.TextAddOrigin, MenuIn)
	}
}

func AddOrigin() func(*tb.Callback) {
	return func(c *tb.Callback) {
		log.Info("BtnAddOrigin clicked")

		log.Debug("message.ID: ", c.Message.ID)

		utils.AddUserState(c.Message.Chat.ID, entity.StateAddOrigin, c.Message.ID+2)

		Bot.Send(c.Sender, entity.TextSendHostMsg)
		Bot.Respond(c, &tb.CallbackResponse{})
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

func ShowMyId() func(*tb.Message) {
	return func(m *tb.Message) {
		log.Info("BtnMyId clicked")

		userChat, message := utils.GetId(m)

		if message != "" {
			Bot.Send(userChat, message)
		}
	}
}

func NewUser() func(*tb.Message) {
	return func(m *tb.Message) {
		log.Info("BtnNewUser clicked")

		MenuIn.Inline(
			MenuIn.Row(BtnAddUser),
		)

		Bot.Send(m.Chat, entity.TextAddUser, MenuIn)
	}
}

func AddUser() func(*tb.Callback) {
	return func(c *tb.Callback) {
		log.Info("BtnAddUser clicked")

		log.Debug("message.ID: ", c.Message.ID)

		utils.AddUserState(c.Message.Chat.ID, entity.StateAddUserEmail, c.Message.ID+2)

		Bot.Send(c.Sender, entity.TextSendEmailMsg)
		Bot.Respond(c, &tb.CallbackResponse{})
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

		log.Debug("send unknown message")
		Bot.Send(m.Chat, entity.TextUnknownMsg)
	}
}
