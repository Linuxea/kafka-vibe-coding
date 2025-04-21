package kafka

import (
	"time"
)

// Config holds the configuration for Kafka client
type Config struct {
	// Brokers is a list of Kafka brokers addresses
	Brokers []string
	// Topic is the Kafka topic
	Topic string
	// GroupID is the consumer group ID
	GroupID string
	// ReadTimeout is the maximum time to wait for a message from Kafka
	ReadTimeout time.Duration
	// WriteTimeout is the maximum time to wait for a message to be written to Kafka
	WriteTimeout time.Duration
	// AutoCreateTopic determines whether topics should be created automatically if they don't exist
	AutoCreateTopic bool
	// NumPartitions defines the number of partitions when auto-creating a topic
	NumPartitions int
	// ReplicationFactor defines the replication factor when auto-creating a topic
	ReplicationFactor int
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		Brokers:           []string{"localhost:9092"},
		Topic:             "test-topic",
		GroupID:           "test-group",
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		AutoCreateTopic:   true,
		NumPartitions:     1,
		ReplicationFactor: 1,
	}
}
