package main

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

func main() {
	// to consume messages
	topic := "test-topic"
	partition := 0

	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:29092", topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	batch := conn.ReadBatch(10e3, 1e6)
	bytes := make([]byte, 10e3)
	for {
		n, err := batch.Read(bytes)
		if err != nil {
			break
		}
		fmt.Println(string(bytes[:n]))
	}

	if err := batch.Close(); err != nil {
		log.Fatal("failed to close batch: ", err)
	}

	if err := conn.Close(); err != nil {
		log.Fatal("failed to close connection:", err)
	}
}
