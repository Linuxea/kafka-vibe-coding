package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/linuxea/kafka-vibe-coding/pkg/kafka"
	"github.com/twmb/franz-go/pkg/kerr"
)

func main() {
	// 创建一个新的消费者配置
	config := kafka.DefaultConsumerConfig()
	config.Topic = "test-topic"
	config.GroupID = "test-consumer-group"
	config.ReadTimeout = 30 * time.Second // 从5秒增加到30秒

	// 创建一个新的消费者
	consumer, err := kafka.NewConsumer(config)
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}
	defer consumer.Close()

	// Create a context with cancellation for clean shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle termination signals
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signals
		fmt.Println("Shutting down consumer...")
		cancel()
	}()

	fmt.Printf("Consumer started. Listening for messages on topic: %s\n", config.Topic)
	fmt.Println("Press Ctrl+C to exit")

	// Consume messages with improved error handling
	for {
		select {
		case <-ctx.Done():
			return
		default:
			// Create a timeout context for each read operation
			readCtx, readCancel := context.WithTimeout(ctx, config.ReadTimeout)
			msg, err := consumer.Consume(readCtx)

			if err != nil {
				readCancel() // Cancel the read context on error

				if err == context.DeadlineExceeded || err == context.Canceled {
					// Just a timeout, not an error
					continue
				}

				// Check if this is a connection error
				if strings.Contains(err.Error(), "connection refused") ||
					strings.Contains(err.Error(), "EOF") {
					log.Printf("Connection lost, waiting to reconnect...")
					time.Sleep(3 * time.Second)
					continue
				}

				// Check if the error is a retriable error
				if isRetriableError(err) {
					log.Printf("Retriable error: %v. Retrying...", err)
					time.Sleep(2 * time.Second)
					continue
				}

				log.Printf("Error consuming message: %v", err)
				time.Sleep(1 * time.Second)
				continue
			}

			// Process the message
			fmt.Printf("Received message: Key=%s, Value=%s\n",
				string(msg.Key), string(msg.Value))

			// Commit the message
			if err := consumer.CommitOffsets(ctx); err != nil {
				log.Printf("Failed to commit offsets: %v", err)
			}

			readCancel() // Always cancel the read context when done
		}
	}
}

// isRetriableError checks if a Kafka error is retriable
func isRetriableError(err error) bool {
	// Check for specific franz-go errors
	if kerr.IsRetriable(err) {
		return true
	}

	// Check for network-related errors and other Kafka errors by error message
	errMsg := err.Error()
	if strings.Contains(errMsg, "context deadline exceeded") ||
		strings.Contains(errMsg, "connection refused") ||
		strings.Contains(errMsg, "broken pipe") ||
		strings.Contains(errMsg, "group") ||
		strings.Contains(errMsg, "member") {
		return true
	}

	return false
}
