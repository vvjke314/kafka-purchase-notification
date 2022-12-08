package kafka

import (
	"context"
	"fmt"
	"github.com/gen2brain/dlgs"
	"github.com/vvjke314/kafka-purchase-notification/telegram"
	"log"
	"os"
	"strconv"
	"time"

	kafkago "github.com/segmentio/kafka-go"
)

type Reader struct {
	Reader *kafkago.Reader
}

func NewKafkaReader() *Reader {
	reader := kafkago.NewReader(kafkago.ReaderConfig{
		Brokers:   []string{"localhost:29092"},
		Topic:     "test-topic",
		Partition: 0,
	})

	return &Reader{
		Reader: reader,
	}
}

func (k *Reader) FetchMessage(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(1 * time.Second):
			off, _ := os.ReadFile("offset.txt")
			ioff, _ := strconv.Atoi(string(off))
			k.Reader.SetOffset(int64(ioff))
			message, err := k.Reader.FetchMessage(ctx)
			if err != nil {
				return err
			}
			var n int
			log.Printf("Напоминание получено и отправлено в канал: %v \n", string(message.Value))
			for n = 0; n < 3; n++ {
				err = telegram.SendMessageService(string(message.Key), string(message.Value))
				if err == nil {
					break
				}
				log.Printf("Не удалось отправить напоминание пользователю, пытаюсь отправить снова...")
				time.Sleep(1 * time.Second)
			}
			f, _ := os.Create("offset.txt")
			if n == 3 {
				k.Reader.SetOffset(k.Reader.Offset() - 1)
				log.Println(k.Reader.Offset())
				log.Printf("Too much tries")
				info := fmt.Sprintf("Не получилось отправить напоминание:{%s} пользователю", string(message.Value))
				_, err := dlgs.Info("Info", info)
				if err != nil {
					panic(err)
				}
			} else {
				k.Reader.SetOffset(k.Reader.Offset())
				log.Println("Offset: ", k.Reader.Offset())
			}
			offsetStr := fmt.Sprintf("%v", k.Reader.Offset())
			f.WriteString(offsetStr)

		}
	}
}
