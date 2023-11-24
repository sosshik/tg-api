package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type HTTPClientInterface interface {
	PostForm(url string, data url.Values) (*http.Response, error)
}

type HandleUserInput interface {
	HandleUserInput(*Api, Update)
}

type Api struct {
	SendString string
	GetUpdStr  string
	HTTPClient HTTPClientInterface
	callback   map[string]func(*Api, Update)
	UserInput  HandleUserInput
	mu         sync.Mutex
}

func (a *Api) AddCallback(f func(*Api, Update), key string) {

	if a.callback == nil {

		a.mu.Lock()

		if a.callback == nil {

			a.callback = make(map[string]func(*Api, Update))

		}

		a.mu.Unlock()

	}

	a.callback[key] = f
	log.Warnf("callback %s was set", key)
}

func ParseTelegramRequest(r *http.Request) (Update, error) {

	var update Update

	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		return Update{}, fmt.Errorf("could not decode incoming update %w", err)
	}

	return update, nil
}

func (a *Api) SendTextToTelegramChat(chatId int, text string) (string, error) {

	response, err := a.HTTPClient.PostForm(a.SendString, url.Values{"chat_id": {strconv.Itoa(chatId)}, "text": {text}})
	if err != nil {
		return "", fmt.Errorf("error when posting text \"%s\" to the chat %d: %w", text, chatId, err)
	}

	defer response.Body.Close()

	switch code := response.StatusCode; {
	case code >= 400 && code < 500:
		return "", fmt.Errorf("bad request: %v", response.Body)
	case code >= 500 && code < 600:
		return "", fmt.Errorf("internal server error")
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("error in parsing telegram response %v, from chat id %d: %w", response, chatId, err)
	}

	return string(body), nil
}

func (a *Api) SendMessage(message OutgoingMessage) error {
	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("error encoding message: %w", err)
	}

	_, err = http.Post(a.SendString, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("error sending message: %w", err)
	}

	return nil
}

func (a *Api) GetUpdates(offset int) ([]Update, error) {
	resp, err := http.Get(a.GetUpdStr + "?offset=" + strconv.Itoa(offset))
	if err != nil {
		return nil, fmt.Errorf("error parsing telegram response: %w", err)
	}
	defer resp.Body.Close()

	var result []Update
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding json: %w", err)
	}

	return result, nil
}

func (a *Api) PollUpdates() {
	offset := 0
	for {
		updates, err := a.GetUpdates(offset)
		if err != nil {
			log.Println("Error getting updates:", err)
			continue
		}

		for _, update := range updates {
			offset = update.UpdateId + 1

			if strings.HasPrefix(update.Message.Text, "/") {
				a.HandleCommand(update)
			} else {
				a.UserInput.HandleUserInput(a, update)
			}
		}

		time.Sleep(1 * time.Second)
	}
}

func (a *Api) HandleCommand(update Update) {
	if f, ok := a.callback[update.Message.Text]; ok {
		f(a, update)
	} else if update.Message.Location.Latitude != 0 && update.Message.Location.Longitude != 0 {

		if f, ok := a.callback["/location"]; ok {

			f(a, update)

		} else {
			log.Warnf("Please add callback for /location command to handle users location")
			a.SendMessageWithLog("*Unknown command*", update)
		}

	} else {
		a.SendMessageWithLog("*Unknown command*", update)
	}
}

func (a *Api) HandleTelegramWebHook(w http.ResponseWriter, r *http.Request) {

	log.Infof("parsing telegram request %v", r)

	update, err := ParseTelegramRequest(r)
	if err != nil {
		log.Warnf("error while parsing update: %v\n", err)
		return
	}

	if strings.HasPrefix(update.Message.Text, "/") {
		a.HandleCommand(update)
	} else {
		a.UserInput.HandleUserInput(a, update)
	}
}

func (a *Api) SendMessageWithLog(text string, update Update) {

	message := OutgoingMessage{
		ChatId:    update.Message.Chat.Id,
		Text:      text,
		ParseMode: "Markdown",
	}

	log.Infof("Sending \" %s\" message to chat_id: %d", text, update.Message.Chat.Id)

	err := a.SendMessage(message)
	if err != nil {
		log.Warnf("got error %s while sending start message to telegram, chat id is %d", err, update.Message.Chat.Id)
		return
	} else {
		log.Infof("message \" %s\" successfuly distributed to chat id %d", text, update.Message.Chat.Id)
	}
}
