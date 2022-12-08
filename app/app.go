package app

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/gen2brain/dlgs"
	"github.com/vvjke314/kafka-purchase-notification/models"
	"io"
	"log"
	"net/http"
)

type App struct {
	ctx context.Context
}

func NewApplication(ctx context.Context) *App {
	return &App{
		ctx: ctx,
	}
}

var Users = map[string]string{
	"Сергей": "719205059",
	"Лев":    "484964911",
}

func (a *App) Run(channel chan models.ResponseMessage) error {
	var names []string
	for name := range Users {
		names = append(names, name)
	}
	user, st, err := dlgs.List("Entry", "Выберите пользователя:", names)
	if err != nil {
		panic(err)
	}
	if user == "" || st == false {
		return errors.New("Canceled")
	}
	note, st, err := dlgs.Entry("Entry", "Введите напоминание:", "")
	if err != nil {
		panic(err)
	}
	if note == "" || st == false {
		return errors.New("Canceled")
	}
	time, st, err := dlgs.Entry("Entry", "Введите время напоминания:", "")
	if err != nil {
		panic(err)
	}
	if time == "" || st == false {
		return errors.New("Canceled")
	}

	resp, err := GetMessage(user, Users[user], time, note)
	message, err := ParseMessage(resp)
	channel <- *message
	if st == false {
		return errors.New("Canceled")
	} else {
		a.Run(channel)
	}
	return err
}

func ParseMessage(resp *http.Response) (*models.ResponseMessage, error) {
	rawMessage, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Can't read response body")
	}

	message := &models.ResponseMessage{}
	err = json.Unmarshal(rawMessage, message)
	if err != nil {
		log.Printf("Can't unmarshall response")
	}
	return message, err
}

func GetMessage(name, id, time, note string) (*http.Response, error) {
	client := http.Client{}
	body := models.RequestMessage{
		Id:   id,
		Name: name,
		Time: time,
		Note: note,
	}
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		log.Printf("Can't marshal req body")
	}
	reqBody := bytes.NewReader(bodyJSON)
	req, err := http.NewRequest("POST", "http://127.0.0.1:8080/notification", reqBody)
	if err != nil {
		log.Printf("Can't send request")
	}
	return client.Do(req)
}
