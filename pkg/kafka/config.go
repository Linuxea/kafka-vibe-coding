package kafka

import (
	"time"
)

// CommonConfig holds the common configuration shared between consumer and producer
type CommonConfig struct {
	// Brokers is a list of Kafka brokers addresses
	Brokers []string
	// Topic is the Kafka topic
	Topic string
}

// ConsumerConfig holds the configuration specific to Kafka consumer
type ConsumerConfig struct {
	// Embedded common configuration
	CommonConfig
	// GroupID is the consumer group ID
	GroupID string
	// ReadTimeout is the maximum time to wait for a message from Kafka
	ReadTimeout time.Duration
}

// ProducerConfig holds the configuration specific to Kafka producer
type ProducerConfig struct {
	// Embedded common configuration
	CommonConfig
	// WriteTimeout is the maximum time to wait for a message to be written to Kafka
	WriteTimeout time.Duration
	// AutoCreateTopic determines whether topics should be created automatically if they don't exist
	AutoCreateTopic bool
	// NumPartitions defines the number of partitions when auto-creating a topic
	NumPartitions int
	// ReplicationFactor defines the replication factor when auto-creating a topic
	ReplicationFactor int
}

// DefaultCommonConfig returns a default common configuration
func DefaultCommonConfig() CommonConfig {
	return CommonConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "test-topic",
	}
}

// DefaultConsumerConfig returns a default consumer configuration
func DefaultConsumerConfig() *ConsumerConfig {
	return &ConsumerConfig{
		CommonConfig: DefaultCommonConfig(),
		GroupID:      "test-group",
		ReadTimeout:  10 * time.Second,
	}
}

// DefaultProducerConfig returns a default producer configuration
func DefaultProducerConfig() *ProducerConfig {
	return &ProducerConfig{
		CommonConfig:      DefaultCommonConfig(),
		WriteTimeout:      10 * time.Second,
		AutoCreateTopic:   true,
		NumPartitions:     1,
		ReplicationFactor: 1,
	}
}
