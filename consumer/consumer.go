package consumer

import (
	"context"
	"eduardopintor/kafka-consumer-microservice.git/config"
	"eduardopintor/kafka-consumer-microservice.git/producer"
	"eduardopintor/kafka-consumer-microservice.git/service"
	"fmt"

	kafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	log "github.com/sirupsen/logrus"
)

type Consumer struct {
	consumer *kafka.Consumer
	topic    string
}

func NewConsumer(topic string) (*Consumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": config.KafkaBootstrapServers,
		"group.id":          config.KafkaGroupId,
		"auto.offset.reset": "earliest",
		"sasl.username":     config.KafkaUsername,
		"sasl.password":     config.KafkaPassword,
		"sasl.mechanism":    config.KafkaSaslMechanism,
		"security.protocol": config.KafkaSecurityProtocol,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka consumer: %v", err)
	}

	return &Consumer{
		consumer: c,
		topic:    topic,
	}, nil
}

func (c *Consumer) Start(ctx context.Context) {
	err := c.consumer.Subscribe(c.topic, nil)
	if err != nil {
		log.Errorf("failed to subscribe to topic %s: %v", c.topic, err)
		return
	}

	for {
		message, err := c.consumer.ReadMessage(-1)
		if err != nil {
			log.Errorf("Error reading message: %v", err)
			continue
		}

		log.Infof("Processing message: %s", string(message.Value))

		orderService, err := service.NewOrder(message)
		if err != nil {
			log.Error("[ProcessMessage] - Error creating order service:", err)
			c.retryProducer(ctx, message)
			return
		}

		response, err := orderService.Process()
		if err != nil {
			log.Error("[ProcessMessage] - Error processing order:", err)
			c.retryProducer(ctx, message)
			return
		}

		log.Infof("Order processed successfully: %s", response)
	}
}

func (c *Consumer) retryProducer(ctx context.Context, message *kafka.Message) error {
	p, err := producer.NewProducer(config.KafkaTopicRetry)
	if err != nil {
		return fmt.Errorf("error creating retry producer: %v", err)
	}

	if err := p.ProduceMessage(ctx, message, 0); err != nil {
		return fmt.Errorf("error producing message to retry topic: %v", err)
	}

	defer p.Close()
	return nil
}

func (c *Consumer) Close() {
	c.consumer.Close()
}
