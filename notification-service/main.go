package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/segmentio/kafka-go"
)

type TaskCreatedEvent struct {
	TaskID      int    `json:"task_id"`
	UserID      int    `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func handleTaskCreatedEvent(data []byte) error {
	var event TaskCreatedEvent

	if err := json.Unmarshal(data, &event); err != nil {
		return err
	}

	log.Printf(
		"Notification: task #%d created for user #%d: %s",
		event.TaskID,
		event.UserID,
		event.Title,
	)

	return nil
}
func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer cancel()

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"kafka:9092"},
		Topic:   "task.created",
		GroupID: "notification-service",
	})

	defer reader.Close()

	log.Println("Notification service started")

	for {
		message, err := reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				log.Println("Notification service stopped")
				return
			}

			log.Println("Kafka error:", err)
			continue
		}

		var event TaskCreatedEvent

		if err := handleTaskCreatedEvent(message.Value); err != nil {
			log.Println("Failed to handle event:", err)
			continue
		}

		log.Printf(
			"Notification: task #%d created for user #%d: %s",
			event.TaskID,
			event.UserID,
			event.Title,
		)
	}
}
