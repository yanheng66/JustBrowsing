package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"github.com/JustBrowsing/query-service/config"
	"go.uber.org/zap"
)

// MessageHandler is a function that processes Kafka messages
type MessageHandler func(topic string, key []byte, value []byte, timestamp time.Time) error

// Consumer represents a Kafka consumer
type Consumer struct {
	consumer sarama.ConsumerGroup
	ready    chan bool
	logger   *zap.Logger
	handler  MessageHandler
}

// ConsumerGroupHandler implements the sarama.ConsumerGroupHandler interface
type ConsumerGroupHandler struct {
	ready   chan bool
	logger  *zap.Logger
	handler MessageHandler
}

// Setup is called when the consumer session is established
func (h *ConsumerGroupHandler) Setup(session sarama.ConsumerGroupSession) error {
	h.logger.Debug("consumer group setup", zap.String("member-id", session.MemberID()))
	close(h.ready)
	return nil
}

// Cleanup is called when the consumer session is terminated
func (h *ConsumerGroupHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	h.logger.Debug("consumer group cleanup", zap.String("member-id", session.MemberID()))
	return nil
}

// ConsumeClaim is called to consume messages from a topic partition
func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		h.logger.Debug("message received",
			zap.String("topic", msg.Topic),
			zap.Int32("partition", msg.Partition),
			zap.Int64("offset", msg.Offset),
			zap.String("key", string(msg.Key)),
			zap.Time("timestamp", msg.Timestamp),
		)

		if err := h.handler(msg.Topic, msg.Key, msg.Value, msg.Timestamp); err != nil {
			h.logger.Error("failed to process message",
				zap.String("topic", msg.Topic),
				zap.Error(err),
			)
		} else {
			// Mark the message as processed
			session.MarkMessage(msg, "")
		}
	}
	return nil
}

// NewConsumer creates a new Kafka consumer
func NewConsumer(cfg config.KafkaConfig, logger *zap.Logger, handler MessageHandler) (*Consumer, error) {
	// Configure the consumer
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Return.Errors = true

	// Create the consumer group
	consumerGroup, err := sarama.NewConsumerGroup(cfg.Brokers, cfg.GroupID, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer group: %w", err)
	}

	consumer := &Consumer{
		consumer: consumerGroup,
		ready:    make(chan bool),
		logger:   logger,
		handler:  handler,
	}

	return consumer, nil
}

// Consume starts consuming messages from the specified topics
func (c *Consumer) Consume(ctx context.Context, topics []string) error {
	// Track errors from the consumer
	go func() {
		for err := range c.consumer.Errors() {
			c.logger.Error("consumer error", zap.Error(err))
		}
	}()

	// Create a new handler for each consumption cycle
	handler := &ConsumerGroupHandler{
		ready:   make(chan bool),
		logger:  c.logger,
		handler: c.handler,
	}

	// Start consuming in a loop
	go func() {
		for {
			// Consume should be run in an infinite loop
			c.logger.Info("starting consumer group", zap.Strings("topics", topics))
			if err := c.consumer.Consume(ctx, topics, handler); err != nil {
				c.logger.Error("consumer group error", zap.Error(err))
			}

			// Check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				c.logger.Info("consumer context cancelled, stopping")
				return
			}

			// Reset the ready channel for the next consumption cycle
			handler.ready = make(chan bool)
		}
	}()

	// Wait until the consumer is set up
	<-handler.ready
	c.logger.Info("consumer is ready", zap.Strings("topics", topics))

	// Keep the consumer running until the context is cancelled
	select {
	case <-ctx.Done():
		c.logger.Info("context cancelled, closing consumer")
		if err := c.consumer.Close(); err != nil {
			return fmt.Errorf("failed to close consumer: %w", err)
		}
	}

	return nil
}