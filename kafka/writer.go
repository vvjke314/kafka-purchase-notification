package kafka

import (
	"context"
	"encoding/json"
	kafkago "github.com/segmentio/kafka-go"
	"github.com/vvjke314/kafka-purchase-notification/models"
	"log"
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

func (k *Writer) WriteMessages(ctx context.Context, responseChannel chan models.ResponseMessage) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case message := <-responseChannel:
			//resp, err := GetMessage()
			//message, err := ParseMessage(resp)
			log.Printf("я тут")
			val, err := json.Marshal(message)
			key, err := json.Marshal(message.Id)
			m := kafkago.Message{
				Key:   key,
				Value: val,
			}
			err = k.Writer.WriteMessages(ctx, m)
			if err != nil {
				return err
			}
			log.Printf("Message %s writed into queue", m.Value)
		}
	}
}

//func ParseMessage(resp *http.Response) (*models.ResponseMessage, error) {
//	rawMessage, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		log.Printf("Can't read response body")
//	}
//
//	message := &models.ResponseMessage{}
//	err = json.Unmarshal(rawMessage, message)
//	if err != nil {
//		log.Printf("Can't unmarshall response")
//	}
//	return message, err
//}
//
//func GetMessage(email string) (*http.Response, error) {
//	product := ds.Products[rand.Intn(5)]
//	client := http.Client{}
//	body := models.RequestMessage{
//		Email:   email,
//		Product: product,
//	}
//	bodyJSON, err := json.Marshal(body)
//	if err != nil {
//		log.Printf("Can't marshal req body")
//	}
//	reqBody := bytes.NewReader(bodyJSON)
//	req, err := http.NewRequest("POST", "http://127.0.0.1:8080/purchase", reqBody)
//	if err != nil {
//		log.Printf("Can't send request")
//	}
//	return client.Do(req)
//}
