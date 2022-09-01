package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"telegram/config"

	log "github.com/sirupsen/logrus"

	tb "gopkg.in/tucnak/telebot.v2"
)

type ResponseData struct {
	Status  int    `json:"status"`
	Origins string `json:"origins"`
}

func GetId(m *tb.Message) (Recipient, string) {
	var userChat Recipient

	if !m.Private() {
		log.Error("Error: chat is not private")
		return Recipient{}, ""
	}

	userChat.ID = int(m.Chat.ID)

	message := "Сообщи этот ID админу для авторизации: " + userChat.Recipient()

	return userChat, message
}

func isAdmin(chatId int) error {
	req, err := setAdminRequest(chatId)
	if err != nil {
		log.Error("Error setting request for admin: ", err)

		return err
	}

	resp, err := BackendClient.Do(&req)
	if err != nil {
		log.Error("Error checking admin from Auth Backend: ", err)

		return err
	}

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("User with ID = %d is not an admin", chatId)
	}

	return nil
}

func setAdminRequest(chatId int) (http.Request, error) {
	body, err := json.Marshal(map[string]string{
		"chat-id": fmt.Sprintf("%d", chatId),
	})
	if err != nil {
		return http.Request{}, fmt.Errorf("Error creating json body: %w", err)
	}

	responseBody := bytes.NewReader(body)

	req, err := http.NewRequest("POST", config.Args.AUTH_URL, responseBody)
	if err != nil {
		return http.Request{}, fmt.Errorf("Error creating request to Backend: %w", err)
	}

	setHeaders(req)

	return *req, nil
}

func setHeaders(req *http.Request) {
	req.Header.Set("X-Green-Origin", "telegram-bot")
	req.Header.Set("Api-Key", config.Args.API_KEY)
}

func getOrigins() (string, error) {
	var data ResponseData

	req, err := http.NewRequest("GET", config.Args.ORIGIN_URL, nil)
	if err != nil {
		return "", fmt.Errorf("Error making GET /origins request: %w", err)
	}

	setHeaders(req)

	resp, err := BackendClient.Do(req)
	if err != nil {
		log.Error("Error calling GET /origins: ", err)

		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Error in response from GET /origins")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Error reading response from GET /origins")
	}

	err = json.Unmarshal([]byte(body), &data)
	if err != nil {
		return "", fmt.Errorf("\nError unmarshalling response JSON %w", err)
	}

	return data.Origins, err
}

func AddUserState(chatID int64, state string, msgID int) {
	if _, userExist := UserStates[chatID]; !userExist {
		UserStates[chatID] = make(map[string]int)
	}

	UserStates[chatID][state] = msgID

	log.Debug("current map: ", UserStates)
}
