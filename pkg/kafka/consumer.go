package kafka

import (
	"context"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
)

// Message represents a Kafka message
type Message struct {
	Key   []byte
	Value []byte
	Time  time.Time
}

// Consumer represents a Kafka consumer
type Consumer struct {
	client *kgo.Client
	config *ConsumerConfig
}

// NewConsumer creates a new Kafka consumer with the given configuration
func NewConsumer(config *ConsumerConfig) (*Consumer, error) {
	opts := []kgo.Opt{
		kgo.SeedBrokers(config.Brokers...),
		kgo.ConsumerGroup(config.GroupID),
		kgo.ConsumeTopics(config.Topic),
		kgo.ClientID("kafka-vibe-consumer"),
		kgo.FetchMaxWait(config.ReadTimeout),
		kgo.DisableAutoCommit(), // We'll handle commits manually for better control
	}

	client, err := kgo.NewClient(opts...)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		client: client,
		config: config,
	}, nil
}

// Consume reads and returns the next message from the consumer
func (c *Consumer) Consume(ctx context.Context) (Message, error) {
	fetches := c.client.PollFetches(ctx)
	if errs := fetches.Errors(); len(errs) > 0 {
		// Return the first error
		return Message{}, errs[0].Err
	}

	// Get the first record from the fetches
	var message Message
	fetches.EachRecord(func(r *kgo.Record) {
		message = Message{
			Key:   r.Key,
			Value: r.Value,
			Time:  r.Timestamp,
		}
		return // Stop iterating after the first record
	})

	return message, nil
}

// ConsumeWithAutoCommit reads messages and automatically commits offsets
func (c *Consumer) ConsumeWithAutoCommit(ctx context.Context) (Message, error) {
	message, err := c.Consume(ctx)
	if err != nil {
		return Message{}, err
	}

	// Commit offsets for this partition
	if err := c.client.CommitUncommittedOffsets(ctx); err != nil {
		return message, err
	}

	return message, nil
}

// Close closes the consumer
func (c *Consumer) Close() error {
	c.client.Close()
	return nil
}

// CommitOffsets commits the current offsets
func (c *Consumer) CommitOffsets(ctx context.Context) error {
	return c.client.CommitUncommittedOffsets(ctx)
}
