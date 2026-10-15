package kafka

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

func StartConsumer(ctx context.Context) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"kafka:9092"},
		Topic:   "task.created",
		GroupID: "taskmanager-group",
	})

	defer reader.Close()

	log.Println("Kafka consumer started")

	for {
		message, err := reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}

			log.Println("Kafka consumer error:", err)
			continue
		}

		log.Printf(
			"Kafka message received: %s\n",
			string(message.Value),
		)
	}
}
