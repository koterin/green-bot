package utils

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"telegram/config"

	log "github.com/sirupsen/logrus"

	tb "gopkg.in/tucnak/telebot.v2"
)

var (
	TextAddUser         = "Тут можно добавить нового пользователя - но обязательно для какого-то определенного сервиса. Например, разрешить пользователю test@test.com доступ к test.example.com"
	TextAddOrigin       = "Тут можно добавить новый сервис для авторизации. Важно! Сам сервис уже должен быть закрыт green-proxy"
	TextAdminRestricted = "Эта опция только для админов."
	TextInternalError   = "Что-то пошло не так. Попробуй еще"
	Menu                = &tb.ReplyMarkup{ResizeReplyKeyboard: true}
	MenuIn              = &tb.ReplyMarkup{}
	BtnNewUser          = Menu.Text("Добавить нового пользователя")
	BtnNewOrigin        = Menu.Text("Подключить новый сервис")
	BtnMyId             = Menu.Text("Мой ID")
	BtnShowOrigins      = MenuIn.Data("Подключенные сервисы", "origins")
	BtnAddOrigin        = MenuIn.Data("Добавить новый", "newOrigin")
	BtnAddUser          = MenuIn.Data("Добавить нового", "newUsers")
	BackendClient       = &http.Client{Timeout: 10 * time.Second}
	Bot                 = &tb.Bot{}
)

type Recipient struct {
	ID int
}

func (user Recipient) Recipient() string {
	return strconv.Itoa(user.ID)
}

func StartTelegramBot(ctx context.Context) {
	settings := tb.Settings{
		Token: config.Args.TG_BOT_KEY,
		Poller: &tb.LongPoller{
			Timeout: 1 * time.Second,
		},
	}

	Bot, _ = tb.NewBot(settings)

	Bot.Handle("/start", OnStart())
	Bot.Handle(&BtnMyId, ShowMyId())
	Bot.Handle(&BtnNewUser, NewUser())
	Bot.Handle(&BtnNewOrigin, NewOrigin())
	Bot.Handle(&BtnShowOrigins, ShowOrigins())

	go func() {
		Bot.Start()
	}()

	log.Info("Telegram Bot started")

	<-ctx.Done()

	log.Info("Telegram Bot stopped")
	Bot.Stop()
}
