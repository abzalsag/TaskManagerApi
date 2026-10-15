package kafka

import (
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer() *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP("kafka:9092"),
			Topic:    "task.created",
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *Producer) PublishTaskCreated(ctx context.Context, event TaskCreatedEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return p.writer.WriteMessages(ctx, kafka.Message{
		Value: data,
	})

}

func (p *Producer) Close() error {
	return p.writer.Close()
}
