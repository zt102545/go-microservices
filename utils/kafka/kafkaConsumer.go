package kafka

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"go-microservices/utils/gos"
	"go-microservices/utils/logs"
	"time"
)

type Consumer struct {
	consumerGroup sarama.ConsumerGroup
	cancel        context.CancelFunc
}

type ConsumerGrp struct {
	topic   []string
	handler func(message *sarama.ConsumerMessage)
}

func (c *ConsumerGrp) Setup(session sarama.ConsumerGroupSession) error {
	fmt.Println("消费者启动")
	return nil
}

func (c *ConsumerGrp) Cleanup(session sarama.ConsumerGroupSession) error {
	fmt.Println("消费者会话结束")
	return nil
}

func (c *ConsumerGrp) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		fmt.Printf("topic=%s, partition=%d, offset=%d, value=%s\n", msg.Topic, msg.Partition, msg.Offset, msg.Value)
		safeHandler(c.handler, msg)
		session.MarkMessage(msg, "")
	}
	return nil
}

func safeHandler(handler func(msg *sarama.ConsumerMessage), msg *sarama.ConsumerMessage) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("recovered from panic in message handler: %v\n", r)
		}
	}()

	handler(msg)
}

func NewConsumerGrp(kafkaAddr []string, groupId string) *Consumer {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Version = sarama.MaxVersion
	config.Consumer.Group.Session.Timeout = 30 * time.Second
	config.Consumer.MaxProcessingTime = 10 * time.Second
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	client, err := sarama.NewClient(kafkaAddr, config)
	if err != nil {
		logs.Err(context.Background(), "kafka client create failed, error: %v", err, logs.Flag("kafka"))
		return nil
	}
	consumerGroup, err := sarama.NewConsumerGroupFromClient(groupId, client)
	if err != nil {
		logs.Err(context.Background(), "kafka consumer create failed, error: %v", err, logs.Flag("kafka"))
		return nil
	}

	instance := &Consumer{
		consumerGroup: consumerGroup,
	}
	return instance
}

func (c *Consumer) MonitorConsume(ctx context.Context, topic []string, handler func(message *sarama.ConsumerMessage)) {
	consumer := &ConsumerGrp{
		topic:   topic,
		handler: handler,
	}
	ctx, cancel := context.WithCancel(ctx)
	c.cancel = cancel
	gos.GoSafe(func() {
		for {
			if ctx.Err() != nil {
				logs.Err(context.Background(), "shutting down, context cancel and consumer close.", logs.Flag("kafka"))
				return
			}

			err := c.consumerGroup.Consume(ctx, consumer.topic, consumer)
			if err != nil {
				fmt.Println("err:", err)
			}
		}
	})
}

func (c *Consumer) Close() error {
	if c.cancel != nil {
		c.cancel()
	}
	if err := c.consumerGroup.Close(); err != nil {
		return fmt.Errorf("failed to close consumer group: %w", err)
	}
	return nil
}
