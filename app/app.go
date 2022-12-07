package app

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/gen2brain/dlgs"
	"github.com/vvjke314/kafka-purchase-notification/ds"
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

var consumers = map[string]ds.User{
	"G321PK": {
		"185404885",
		false,
	},
	"G444PP": {
		"185404885",
		false,
	},
}

func (a *App) Run(channel chan models.ResponseMessage) error {
	var plates []string
	for pl := range consumers {
		plates = append(plates, pl)
	}
	plate, st, err := dlgs.List("List", "Выберете номер машины:", plates)
	if err != nil {
		panic(err)
	}
	if plate == "" || st == false {
		return errors.New("Canceled")
	}
	var choice []string
	if consumers[plate].IsEnter {
		choice = append(choice, "Вы выехали с парковки!")
		consumers[plate] = ds.User{
			consumers[plate].Id,
			!consumers[plate].IsEnter,
		}
	} else {
		choice = append(choice, "Вы въехали на парковку!")
		consumers[plate] = ds.User{
			consumers[plate].Id,
			!consumers[plate].IsEnter,
		}
	}
	state, st, err := dlgs.List("List", "Select item from list:", choice)
	if err != nil {
		panic(err)
	}
	if state == "" || st == false {
		return errors.New("Canceled")
	}
	resp, err := GetMessage(plate, consumers[plate].Id, state)
	message, err := ParseMessage(resp)
	channel <- *message
	if err != nil {
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

func GetMessage(plate, id, status string) (*http.Response, error) {
	client := http.Client{}
	body := models.RequestMessage{
		Plate:  plate,
		Id:     id,
		Status: status,
	}
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		log.Printf("Can't marshal req body")
	}
	reqBody := bytes.NewReader(bodyJSON)
	req, err := http.NewRequest("POST", "http://127.0.0.1:8080/purchase", reqBody)
	if err != nil {
		log.Printf("Can't send request")
	}
	return client.Do(req)
}
