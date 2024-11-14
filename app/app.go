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

func (a *App) Run(channel chan models.ResponseMessage) error {
	email, st, err := dlgs.Entry("Entry", "Enter your email:", "volohajks@inbox.ru")
	if err != nil {
		panic(err)
	}
	if email == "" || st == false {
		return errors.New("Canceled")
	}
	product, st, err := dlgs.Entry("Entry", "Enter product name:", "")
	resp, err := GetMessage(email, product)
	message, err := ParseMessage(resp)
	channel <- *message
	if product == "" || email == "" || st == false {
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

func GetMessage(email, product string) (*http.Response, error) {
	client := http.Client{}
	body := models.RequestMessage{
		Email:   email,
		Product: product,
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
