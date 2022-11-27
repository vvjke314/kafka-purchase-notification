package kafka

import (
	"bytes"
	"context"
	"encoding/json"
	kafkago "github.com/segmentio/kafka-go"
	"github.com/vvjke314/kafka-purchase-notification/ds"
	"github.com/vvjke314/kafka-purchase-notification/models"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type Writer struct {
	Writer *kafkago.Writer
}

func NewKafkaWriter() *Writer {
	writer := &kafkago.Writer{
		Addr:  kafkago.TCP("localhost:29092"),
		Topic: "test-topic",
	}
	return &Writer{
		Writer: writer,
	}
}

func (k *Writer) WriteMessages(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(1 * time.Second):
			resp, err := GetMessage()
			message, err := ParseMessage(resp)
			val, err := json.Marshal(message)
			key, err := json.Marshal(message.Email)
			m := kafkago.Message{
				Key:   key,
				Value: val,
			}
			err = k.Writer.WriteMessages(ctx, m)
			if err != nil {
				return err
			}
		}
	}
}

func ParseMessage(resp *http.Response) (*models.ResponseMessage, error) {
	rawMessage, err := ioutil.ReadAll(resp.Body)
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

func GetMessage() (*http.Response, error) {
	email := ds.Emails[rand.Intn(2)]
	product := ds.Products[rand.Intn(3)]
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
