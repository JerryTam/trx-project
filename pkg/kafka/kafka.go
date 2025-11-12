package kafka

import (
	"context"
	"fmt"
	"trx-project/pkg/config"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Producer struct {
	writer *kafka.Writer
	logger *zap.Logger
}

type Consumer struct {
	reader *kafka.Reader
	logger *zap.Logger
}

// NewProducer creates a new Kafka producer
func NewProducer(cfg *config.KafkaConfig, logger *zap.Logger) *Producer {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(cfg.Brokers...),
		Balancer: &kafka.LeastBytes{},
	}

	logger.Info("Kafka producer initialized")
	return &Producer{
		writer: writer,
		logger: logger,
	}
}

// SendMessage sends a message to Kafka topic
func (p *Producer) SendMessage(ctx context.Context, topic, key string, value []byte) error {
	msg := kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: value,
	}

	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		p.logger.Error("Failed to send message to Kafka",
			zap.String("topic", topic),
			zap.Error(err))
		return fmt.Errorf("failed to send message: %w", err)
	}

	p.logger.Debug("Message sent to Kafka",
		zap.String("topic", topic),
		zap.String("key", key))
	return nil
}

// Close closes the Kafka producer
func (p *Producer) Close() error {
	return p.writer.Close()
}

// NewConsumer creates a new Kafka consumer
func NewConsumer(cfg *config.KafkaConfig, topic string, logger *zap.Logger) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  cfg.Brokers,
		GroupID:  cfg.GroupID,
		Topic:    topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	logger.Info("Kafka consumer initialized", zap.String("topic", topic))
	return &Consumer{
		reader: reader,
		logger: logger,
	}
}

// ReadMessage reads a message from Kafka
func (c *Consumer) ReadMessage(ctx context.Context) (kafka.Message, error) {
	msg, err := c.reader.ReadMessage(ctx)
	if err != nil {
		c.logger.Error("Failed to read message from Kafka", zap.Error(err))
		return kafka.Message{}, fmt.Errorf("failed to read message: %w", err)
	}

	c.logger.Debug("Message received from Kafka",
		zap.String("topic", msg.Topic),
		zap.Int("partition", msg.Partition),
		zap.Int64("offset", msg.Offset))
	return msg, nil
}

// Close closes the Kafka consumer
func (c *Consumer) Close() error {
	return c.reader.Close()
}
