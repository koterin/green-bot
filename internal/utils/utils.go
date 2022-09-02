package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"telegram/config"
	"telegram/internal/entity"

	log "github.com/sirupsen/logrus"
	tb "gopkg.in/telebot.v3"
)

var (
	BackendClient = &http.Client{Timeout: 10 * time.Second}
	UserStates    = make(map[int64]map[string]int) // map['chatID'] = map'btnAddOrigin' = 'message.ID'
)

func GetId(m *tb.Message) string {
	var userChat entity.Recipient

	if !m.Private() {
		log.Error("Error: chat is not private")
		return ""
	}

	userChat.ID = int(m.Chat.ID)

	message := "Сообщи этот ID админу для авторизации: " + userChat.Recipient()

	return message
}

func IsAdmin(chatId int) error {
	req, err := setRequest(map[string]string{
		"chat-id": fmt.Sprintf("%d", chatId),
	}, config.Args.AUTH_URL)
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

func setRequest(payload map[string]string, url string) (http.Request, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return http.Request{}, fmt.Errorf("Error creating json body: %w", err)
	}

	responseBody := bytes.NewReader(body)

	req, err := http.NewRequest("POST", url, responseBody)
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

func GetOrigins() (string, error) {
	var data entity.ResponseData

	req, err := http.NewRequest("GET", config.Args.ORIGIN_URL, nil)
	if err != nil {
		return "", fmt.Errorf("Error making GET /origins request: %w", err)
	}

	setHeaders(req)

	resp, err := BackendClient.Do(req)
	if err != nil {
		log.Error("\nError calling GET /origins: ", err)

		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("\nError in response from GET /origins")
	}

	err = readJson(resp.Body, &data)
	if err != nil {
		return "", fmt.Errorf("\nError in readJson: %w", err)
	}

	return data.Origins, err
}

func readJson(resp io.ReadCloser, data *entity.ResponseData) error {
	body, err := io.ReadAll(resp)
	if err != nil {
		return fmt.Errorf("\nError reading response from GET /origins")
	}

	err = json.Unmarshal([]byte(body), data)
	if err != nil {
		return fmt.Errorf("\nError unmarshalling response JSON %w", err)
	}

	return err
}

func AddUserState(chatID int64, state string, msgID int) {
	if _, userExist := UserStates[chatID]; !userExist {
		UserStates[chatID] = make(map[string]int)
	}

	UserStates[chatID][state] = msgID

	log.Debug("current map: ", UserStates)
}

func ValidateOrigin(origin string) string {
	var data entity.ResponseData

	req, err := setRequest(map[string]string{
		"origin": origin,
	}, config.Args.NEW_ORIGIN_URL)
	if err != nil {
		log.Error("Error setting request for .ValidateOrigin: ", err)

		return entity.TextInternalError
	}

	resp, err := BackendClient.Do(&req)
	if err != nil {
		log.Error("Error creating new Origin: ", err)

		return entity.TextInternalError
	}

	if resp.StatusCode == http.StatusInternalServerError ||
		resp.StatusCode == http.StatusBadRequest {
		return entity.TextInternalError
	}

	err = readJson(resp.Body, &data)
	if err != nil {
		log.Error("Error in .ValidateOrigin: ", err)

		return entity.TextInternalError
	}

	return data.Response
}
