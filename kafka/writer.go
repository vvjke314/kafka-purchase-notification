package kafka

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	kafkago "github.com/segmentio/kafka-go"
	"log"
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
		m := kafkago.Message{
			Value: []byte(fmt.Sprint(uuid.New())),
		}
		err := k.Writer.WriteMessages(ctx, m)
		if err != nil {
			return err
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(4):
			log.Println("Записываем сообщение " + string(m.Value))
		}
	}
}
