package kafka

import (
	"context"
	"github.com/pkg/errors"
	"github.com/vvjke314/kafka-purchase-notification/mail"
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
		case <-time.After(5 * time.Second):
			message, err := k.Reader.FetchMessage(ctx)
			if err != nil {
				return err
			}
			var n int
			log.Printf("message fetched and sent to a channel: %v \n", string(message.Value))
			for n = 0; n < 3; n++ {
				err = mail.SendMessageService(string(message.Key), message.Value)
				if err == nil {
					break
				}
				log.Printf("Can't send message to user, retrying...")
			}
			if n == 3 {
				log.Printf("Too much tries")
				k.Reader.SetOffset(k.Reader.Offset() - 1)
			} else {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case messageCommitChan <- message:
				}
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
