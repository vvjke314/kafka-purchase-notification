package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gen2brain/dlgs"
	"github.com/vvjke314/kafka-purchase-notification/models"
	"github.com/vvjke314/kafka-purchase-notification/vk"
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
		Brokers:   []string{"localhost:29092"}, //надо занести в файлы конфигурации
		Topic:     "test-topic",
		Partition: 0,
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
		case <-time.After(1 * time.Second):
			off, _ := os.ReadFile("offset.txt")
			ioff, _ := strconv.Atoi(string(off))
			k.Reader.SetOffset(int64(ioff))
			message, err := k.Reader.FetchMessage(ctx)
			if err != nil {
				return err
			}
			log.Println(string(message.Value))
			var n int
			VkID, err := strconv.Atoi(string(message.Key))
			if err != nil {
				return err
			}
			resp := models.ResponseMessage{}
			json.Unmarshal(message.Value, &resp)
			log.Printf("message fetched and sent to a channel: %v \n", resp.Status)
			for n = 0; n < 3; n++ {
				err = vk.SendMessageService(VkID, resp.Status)
				if err == nil {
					break
				}
				log.Printf("Can't send message to user, retrying...")
				time.Sleep(1 * time.Second)
			}
			f, _ := os.Create("offset.txt")
			defer f.Close()
			if n == 3 {
				k.Reader.SetOffset(k.Reader.Offset() - 1)
				log.Println(k.Reader.Offset())
				log.Printf("Too much tries")
				info := fmt.Sprintf("Can't send message to user:%s. Call the number:%s", message.Key, resp.Number)
				_, err := dlgs.Info("Info", info)
				if err != nil {
					log.Println(err)
				}
			} else {
				k.Reader.SetOffset(k.Reader.Offset())
				log.Println(k.Reader.Offset())
			}
			offsetStr := fmt.Sprintf("%v", k.Reader.Offset())
			f.WriteString(offsetStr)
		}
	}
}
