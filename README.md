# Kafka-Vibe-Coding

一个基于 Go 语言的高性能 Kafka 客户端库，提供简洁易用的 API 来与 Apache Kafka 进行交互。

## 特性

- 简洁的 API 设计，易于使用
- 支持生产者和消费者模式
- 基于 [franz-go](https://github.com/twmb/franz-go) Kafka 客户端库
- 支持自动创建主题
- 提供同步和异步生产消息的能力
- 灵活的配置选项
- 完整的错误处理机制

## 安装

使用 Go modules 安装:

```bash
go get github.com/linuxea/kafka-vibe-coding
```

## 快速开始

### 生产者示例

```go
package main

import (
    "context"
    "log"
    "time"

    "github.com/linuxea/kafka-vibe-coding/pkg/kafka"
)

func main() {
    // 创建生产者配置
    config := kafka.DefaultProducerConfig()
    config.Topic = "my-topic"

    // 创建生产者
    producer, err := kafka.NewProducer(config)
    if err != nil {
        log.Fatalf("Failed to create producer: %v", err)
    }
    defer producer.Close()

    // 同步发送消息
    ctx := context.Background()
    key := []byte("key")
    value := []byte("Hello Kafka!")
    
    if err := producer.Produce(ctx, key, value); err != nil {
        log.Fatalf("Failed to produce message: %v", err)
    }
    
    log.Println("Message sent successfully!")
}
```

### 消费者示例

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "os/signal"
    "syscall"

    "github.com/linuxea/kafka-vibe-coding/pkg/kafka"
)

func main() {
    // 创建消费者配置
    config := kafka.DefaultConsumerConfig()
    config.Topic = "my-topic"
    config.GroupID = "my-consumer-group"

    // 创建消费者
    consumer, err := kafka.NewConsumer(config)
    if err != nil {
        log.Fatalf("Failed to create consumer: %v", err)
    }
    defer consumer.Close()

    // 处理终止信号
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    signals := make(chan os.Signal, 1)
    signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

    go func() {
        <-signals
        cancel()
    }()

    // 消费消息
    for {
        select {
        case <-ctx.Done():
            return
        default:
            msg, err := consumer.Consume(ctx)
            if err != nil {
                log.Printf("Error consuming message: %v", err)
                continue
            }

            fmt.Printf("Received message: Key=%s, Value=%s\n",
                string(msg.Key), string(msg.Value))

            // 提交偏移量
            if err := consumer.CommitOffsets(ctx); err != nil {
                log.Printf("Failed to commit offsets: %v", err)
            }
        }
    }
}
```

## 详细文档

### 配置选项

#### 通用配置

```go
type CommonConfig struct {
    // Brokers 是 Kafka broker 地址列表
    Brokers []string
    // Topic 是 Kafka 主题
    Topic string
}
```

#### 生产者配置

```go
type ProducerConfig struct {
    // 内嵌通用配置
    CommonConfig
    // WriteTimeout 是等待消息写入 Kafka 的最大时间
    WriteTimeout time.Duration
    // AutoCreateTopic 决定是否在主题不存在时自动创建
    AutoCreateTopic bool
    // NumPartitions 定义自动创建主题时的分区数
    NumPartitions int
    // ReplicationFactor 定义自动创建主题时的复制因子
    ReplicationFactor int
}
```

#### 消费者配置

```go
type ConsumerConfig struct {
    // 内嵌通用配置
    CommonConfig
    // GroupID 是消费者组 ID
    GroupID string
    // ReadTimeout 是从 Kafka 等待消息的最大时间
    ReadTimeout time.Duration
}
```

### 生产者 API

#### 创建生产者

```go
producer, err := kafka.NewProducer(config)
```

#### 同步发送消息

```go
err := producer.Produce(ctx, key, value)
```

#### 异步发送消息

```go
producer.ProduceAsync(key, value, func(record *kgo.Record, err error) {
    // 处理回调
})
```

#### 批量发送消息

```go
messages := []kafka.Message{
    {Key: []byte("key1"), Value: []byte("value1")},
    {Key: []byte("key2"), Value: []byte("value2")},
}
err := producer.ProduceMessages(ctx, messages)
```

### 消费者 API

#### 创建消费者

```go
consumer, err := kafka.NewConsumer(config)
```

#### 消费消息

```go
msg, err := consumer.Consume(ctx)
```

#### 自动提交偏移量的消费

```go
msg, err := consumer.ConsumeWithAutoCommit(ctx)
```

#### 手动提交偏移量

```go
err := consumer.CommitOffsets(ctx)
```

## 示例

完整的生产者和消费者示例可以在 `examples` 目录中找到:

- [生产者示例](examples/producer/main.go)
- [消费者示例](examples/consumer/main.go)

## 依赖项

- [github.com/twmb/franz-go](https://github.com/twmb/franz-go): 高性能的 Go Kafka 客户端
- [github.com/twmb/franz-go/pkg/kadm](https://github.com/twmb/franz-go/tree/master/pkg/kadm): Kafka 管理工具

## 贡献

欢迎提交 Pull Request 和 Issue!

## 许可证

MIT