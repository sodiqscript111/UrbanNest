package kafka

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"time"
)

type Producer struct {
	Writer *kafka.Writer
}

func NewProducer(brokers []string, topic string) *Producer {
	return &Producer{
		Writer: kafka.NewWriter(kafka.WriterConfig{
			Brokers:  brokers,
			Balancer: &kafka.LeastBytes{},
			Topic:    topic,
		}),
	}
}

func (p *Producer) PublishMessage(ctx context.Context, key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return p.Writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(key),
		Value: data,
		Time:  time.Now(),
	})
}

func (p *Producer) Close() error {
	return p.Writer.Close()
}
