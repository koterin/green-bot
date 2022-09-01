package utils

import (
	log "github.com/sirupsen/logrus"

	tb "gopkg.in/tucnak/telebot.v2"
)

var (
	TextAddUser         = "Тут можно добавить нового пользователя - но обязательно для какого-то определенного сервиса. Например, разрешить пользователю test@test.com доступ к test.example.com"
	TextAddOrigin       = "Тут можно добавить новый сервис для авторизации. Важно! Сам сервис уже должен быть закрыт green-proxy"
	TextAdminRestricted = "Эта опция только для админов."
	TextInternalError   = "Что-то пошло не так. Попробуй еще"
	TextSendOriginName  = "Отправьте название нового сервиса:"
	Menu                = &tb.ReplyMarkup{ResizeReplyKeyboard: true}
	MenuIn              = &tb.ReplyMarkup{}
	BtnNewUser          = Menu.Text("Добавить нового пользователя")
	BtnNewOrigin        = Menu.Text("Подключить новый сервис")
	BtnMyId             = Menu.Text("Мой ID")
	BtnShowOrigins      = MenuIn.Data("Подключенные сервисы", "origins")
	BtnAddOrigin        = MenuIn.Data("Добавить новый", "newOrigin")
	BtnAddUser          = MenuIn.Data("Добавить нового", "newUser")
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

func AddOrigin() func(*tb.Callback) {
	return func(c *tb.Callback) {
		log.Info("BtnAddOrigin clicked")

		log.Debug("message.ID: ", c.Message.ID)

		AddUserState(c.Message.Chat.ID, "btnAddOrigin", c.Message.ID+2)

		Bot.Send(c.Sender, "Введите хост нового сервиса:")
		Bot.Respond(c, &tb.CallbackResponse{})
	}
}

func ShowOrigins() func(*tb.Callback) {
	return func(c *tb.Callback) {
		log.Info("BtnShowOrigins clicked")

		origins, err := getOrigins()
		if err != nil {
			log.Info(err)
			Bot.Send(c.Sender, TextInternalError)

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

func AddUser() func(*tb.Callback) {
	return func(c *tb.Callback) {
		log.Info("BtnAddUser clicked")

		log.Debug("message.ID: ", c.Message.ID)

		AddUserState(c.Message.Chat.ID, "btnAddUser", c.Message.ID+2)

		Bot.Send(c.Sender, "Введите почту нового пользователя:")
		Bot.Respond(c, &tb.CallbackResponse{})
	}
}

func OnText() func(*tb.Message) {
	return func(m *tb.Message) {
		var (
			msgWanted bool
			state     string
			msgID     int
		)

		log.Info("Text received")
		log.Debug("m.Chat.ID: ", m.Chat.ID)
		log.Debug("m.ID: ", m.ID)

		if _, userExist := UserStates[m.Chat.ID]; !userExist {
			UserStates[m.Chat.ID] = make(map[string]int)
		}

		for state, msgID = range UserStates[m.Chat.ID] {
			if msgID == m.ID {
				msgWanted = true

				break
			}
		}

		if !msgWanted {
			return
		}

		text := "I wanted this message for state = " + state
		Bot.Send(m.Chat, text)
	}
}
