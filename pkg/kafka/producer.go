package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
)

// Producer represents a Kafka producer
type Producer struct {
	client *kgo.Client
	config *ProducerConfig
}

// NewProducer creates a new Kafka producer with the given configuration
func NewProducer(config *ProducerConfig) (*Producer, error) {
	// Set up producer options
	opts := []kgo.Opt{
		kgo.SeedBrokers(config.Brokers...),
		kgo.ClientID("kafka-vibe-producer"),
		kgo.ProducerBatchMaxBytes(1024 * 1024), // 1MB
		kgo.ProducerLinger(time.Millisecond * 10),
		kgo.RequestTimeoutOverhead(config.WriteTimeout),
		kgo.AllowAutoTopicCreation(), // Auto-create topics if needed
	}

	client, err := kgo.NewClient(opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka client: %w", err)
	}

	producer := &Producer{
		client: client,
		config: config,
	}

	// Check if topic needs to be created with specific settings
	if config.AutoCreateTopic {
		if err := producer.createTopicIfNotExists(config); err != nil {
			// Log the error but continue - client can still try to produce
			fmt.Printf("Warning: Failed to create topic %s: %v\n", config.Topic, err)
		}
	}

	return producer, nil
}

// createTopicIfNotExists checks if a topic exists and creates it if it doesn't
func (p *Producer) createTopicIfNotExists(config *ProducerConfig) error {
	// Create an admin client
	adminClient := kadm.NewClient(p.client)

	// Try to get topic metadata to check if it exists
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	topics, err := adminClient.ListTopics(ctx)
	if err != nil {
		return fmt.Errorf("failed to list topics: %w", err)
	}

	// Check if our topic exists
	topicExists := false
	for _, topic := range topics {
		if topic.Topic == config.Topic {
			topicExists = true
			break
		}
	}

	if !topicExists {
		// Topic doesn't exist, create it with specific settings
		resp, err := adminClient.CreateTopics(ctx, int32(config.NumPartitions), int16(config.ReplicationFactor), nil, config.Topic)
		if err != nil {
			return fmt.Errorf("failed to create topic: %w", err)
		}

		// Check for errors in the response
		for _, topicResult := range resp {
			if topicResult.Err != nil {
				return fmt.Errorf("failed to create topic %s: %v", topicResult.Topic, topicResult.Err)
			}
		}
	}

	return nil
}

// Produce sends a message to Kafka
func (p *Producer) Produce(ctx context.Context, key, value []byte) error {
	record := &kgo.Record{
		Topic: p.config.Topic,
		Key:   key,
		Value: value,
	}

	err := p.client.ProduceSync(ctx, record).FirstErr()
	if err != nil {
		return fmt.Errorf("failed to produce message: %w", err)
	}

	return nil
}

// ProduceAsync sends a message to Kafka asynchronously
func (p *Producer) ProduceAsync(key, value []byte, callback func(*kgo.Record, error)) {
	record := &kgo.Record{
		Topic: p.config.Topic,
		Key:   key,
		Value: value,
	}

	p.client.Produce(context.Background(), record, callback)
}

// ProduceMessages sends multiple messages to Kafka
func (p *Producer) ProduceMessages(ctx context.Context, messages []Message) error {
	// Convert our Message type to kgo.Record
	records := make([]*kgo.Record, len(messages))
	for i, msg := range messages {
		records[i] = &kgo.Record{
			Topic: p.config.Topic,
			Key:   msg.Key,
			Value: msg.Value,
		}
	}

	// Use the producer to send all messages
	for _, record := range records {
		err := p.client.ProduceSync(ctx, record).FirstErr()
		if err != nil {
			return fmt.Errorf("failed to produce message: %w", err)
		}
	}

	return nil
}

// Close closes the producer
func (p *Producer) Close() error {
	p.client.Close()
	return nil
}
