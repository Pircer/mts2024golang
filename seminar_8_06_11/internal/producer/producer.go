package producer

import (
	"context"
	"encoding/json"
	"github.com/twmb/franz-go/pkg/kgo"
	"mts2024golang/seminar_8_06_11/internal/models"
)

type Producer struct {
	client *kgo.Client
	topic  string
}

func New(brokers []string, topic string) (*Producer, error) {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
	)
	if err != nil {
		return nil, err
	}
	return &Producer{client: client, topic: topic}, nil
}
func (p *Producer) SendMessage(user, message string) {
	ctx := context.Background()
	msg := models.Message{User: user, Message: message}
	b, _ := json.Marshal(msg)
	p.client.Produce(ctx, &kgo.Record{Topic: p.topic, Value: b}, nil)
}
func (p *Producer) Close() {
	p.client.Close()
}
