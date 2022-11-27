package main

import (
	"context"
	kafkago "github.com/segmentio/kafka-go"
	"log"

	"github.com/vvjke314/kafka-purchase-notification/kafka"
	"golang.org/x/sync/errgroup"
)

func main() {
	reader := kafka.NewKafkaReader()
	writer := kafka.NewKafkaWriter()

	ctx := context.Background()
	commitMessage := make(chan kafkago.Message)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return reader.FetchMessage(ctx, commitMessage)
	})

	g.Go(func() error {
		return writer.WriteMessages(ctx)
	})

	g.Go(func() error {
		return reader.CommitMessages(ctx, commitMessage)
	})

	err := g.Wait()
	if err != nil {
		log.Fatalln(err)
	}
}
