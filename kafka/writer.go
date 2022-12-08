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
			log.Printf("Напоминание %s записано в очередь", m.Value)
		}
	}
}
