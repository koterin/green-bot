package utils

import (
	log "github.com/sirupsen/logrus"

	tb "gopkg.in/tucnak/telebot.v2"
)

func OnStart() func(*tb.Message) {
	return func(m *tb.Message) {
		userChat, message := GetId(m)

		if message != "" {
			if err := isAdmin(userChat.ID); err != nil {
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

		Bot.Send(m.Chat, TextAddOrigin, MenuIn)
	}
}

func ShowOrigins() func(*tb.Message) {
	return func(m *tb.Message) {
		log.Info("BtnShowOrigins clicked")

		userChat, message := GetId(m)
		if message != "" {
			if err := isAdmin(userChat.ID); err != nil {
				log.Info(err)
				Bot.Send(m.Chat, TextAdminRestricted)

				return
			}
		}

		origins, err := getOrigins()
		if err != nil {
			log.Info(err)
			Bot.Send(m.Chat, TextInternalError)

			return
		}

		Bot.Send(m.Chat, origins, MenuIn)
	}
}

func ShowMyId() func(*tb.Message) {
	return func(m *tb.Message) {
		log.Info("BtnMyId clicked")

		userChat, message := GetId(m)

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

		Bot.Send(m.Chat, TextAddUser, MenuIn)
	}
}
