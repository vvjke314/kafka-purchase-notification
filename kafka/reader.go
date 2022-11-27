package kafka

import (
	"context"
	"github.com/pkg/errors"
	"log"
	"time"

	kafkago "github.com/segmentio/kafka-go"
)

type Reader struct {
	Reader *kafkago.Reader
}

func NewKafkaReader() *Reader {
	reader := kafkago.NewReader(kafkago.ReaderConfig{
		Brokers: []string{"localhost:29092"}, //надо занести в файлы конфигурации
		Topic:   "test-topic",
		GroupID: "group",
	})

	return &Reader{
		Reader: reader,
	}
}

func (k *Reader) FetchMessage(ctx context.Context, messageCommitChan chan kafkago.Message) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(2):
			message, err := k.Reader.FetchMessage(ctx)
			if err != nil {
				return err
			}
			log.Printf("message fetched and sent to a channel: %v \n", string(message.Value))
			select {
			case <-ctx.Done():
				return ctx.Err()
			case messageCommitChan <- message:
			}
		}
	}
}

func (k *Reader) CommitMessages(ctx context.Context, messageCommitChan chan kafkago.Message) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg := <-messageCommitChan:
			err := k.Reader.CommitMessages(ctx, msg)
			if err != nil {
				return errors.Wrap(err, "Reader.CommitMessages")
			}
			log.Printf("committed an msg: %v \n", string(msg.Value))
		}
	}
}
