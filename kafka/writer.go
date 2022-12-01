package kafka

import (
	"context"
	"encoding/json"
	kafkago "github.com/segmentio/kafka-go"
	"github.com/vvjke314/kafka-purchase-notification/models"
	"io/ioutil"
	"log"
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
		case <-time.After(4 * time.Second):
			resp, err := GetMessage()
			message, err := ParseMessage(resp)
			log.Println(message.Status)
			val, err := json.Marshal(message.Status)
			key, err := json.Marshal(message.ID)
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
	req, err := http.Get("http://0.0.0.0:8080/delivery")
	if err != nil {
		log.Printf("Can't send request")
	}
	return req, err
}
