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
	"strconv"
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
	ID, st, err := dlgs.Entry("Entry", "Enter your vkID:", "")
	if err != nil {
		panic(err)
	}
	if ID == "" || st == false {
		return errors.New("Canceled")
	}
	vkID, err := strconv.Atoi(ID)
	if err != nil {
		return errors.New("VkID not int")
	}
	status, st, err := dlgs.Entry("Entry", "Enter status order:", "")
	if status == "" || st == false {
		return errors.New("Canceled")
	}
	resp, err := GetMessage(vkID, status)
	message, err := ParseMessage(resp)
	channel <- *message
	a.Run(channel)
	return nil
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

func GetMessage(ID int, status string) (*http.Response, error) {
	client := http.Client{}
	body := models.RequestMessage{
		ID:     ID,
		Status: status,
	}
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		log.Printf("Can't marshal req body")
	}
	reqBody := bytes.NewReader(bodyJSON)
	req, err := http.NewRequest("POST", "http://127.0.0.1:8080/delivery", reqBody)
	if err != nil {
		log.Printf("Can't send request")
	}
	return client.Do(req)
}
