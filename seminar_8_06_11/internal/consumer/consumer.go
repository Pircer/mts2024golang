package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/twmb/franz-go/pkg/kgo"
	"mts2024golang/seminar_8_06_11/internal/models"
)

type Consumer struct {
	client *kgo.Client
	topic  string
}

func New(brokers []string, topic string) (*Consumer, error) {
	groupID := uuid.New().String()
	client, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
		kgo.ConsumerGroup(groupID),
		kgo.ConsumeTopics(topic),
		kgo.ConsumeResetOffset(kgo.NewOffset().AtStart()),
	)
	if err != nil {
		return nil, err
	}
	return &Consumer{client: client, topic: topic}, nil
}

func (c *Consumer) PrintMessages() {
	ctx := context.Background()
	for {
		fetches := c.client.PollFetches(ctx)
		iter := fetches.RecordIter()
		for !iter.Done() {
			record := iter.Next()
			var msg models.Message
			if err := json.Unmarshal(record.Value, &msg); err != nil {
				fmt.Printf("Error decoding message: %v\n", err)
				continue
			}
			fmt.Printf("%s: %s\n", msg.User, msg.Message)
		}
	}
}

func (c *Consumer) Close() {
	c.client.Close()
}
