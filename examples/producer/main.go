package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/linuxea/kafka-vibe-coding/pkg/kafka"
)

func main() {
	// Create a new configuration with custom settings
	config := &kafka.Config{
		Brokers:           []string{"localhost:9092"},
		Topic:             "test-topic",
		WriteTimeout:      5 * time.Second,
		AutoCreateTopic:   true,
		NumPartitions:     1,
		ReplicationFactor: 1,
	}

	// Create a new producer
	producer, err := kafka.NewProducer(config)
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}
	defer producer.Close()

	// Create a context with cancellation for clean shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle termination signals
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signals
		fmt.Println("Shutting down producer...")
		cancel()
	}()

	// Produce messages until context is cancelled
	counter := 0
	for {
		select {
		case <-ctx.Done():
			return
		default:
			counter++
			key := []byte("现在是北京时间 " + time.Now().Format("2006-01-02 15:04:05") + " key-" + strconv.Itoa(counter))
			value := []byte(fmt.Sprintf("message #%d at %v", counter, time.Now().Format(time.RFC3339)))

			fmt.Printf("Producing message: %s\n", value)
			if err := producer.Produce(ctx, key, value); err != nil {
				log.Printf("Failed to produce message: %v", err)
				time.Sleep(1 * time.Second)
				continue
			}

			time.Sleep(1 * time.Second)
		}
	}
}
