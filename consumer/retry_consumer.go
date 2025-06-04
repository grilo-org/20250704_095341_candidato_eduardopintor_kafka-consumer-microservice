package consumer

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"eduardopintor/kafka-consumer-microservice.git/config"
	"eduardopintor/kafka-consumer-microservice.git/producer"
	"eduardopintor/kafka-consumer-microservice.git/service"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	log "github.com/sirupsen/logrus"
)

type RetryConsumer struct {
	consumer     *kafka.Consumer
	producer     *producer.Producer
	topic        string
	maxRetries   int
	delayMinutes int
}

func NewRetryConsumer(topic string) (*RetryConsumer, error) {
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
		return nil, err
	}

	p, err := producer.NewProducer(config.KafkaTopicRetry)
	if err != nil {
		return nil, fmt.Errorf("error creating producer: %v", err)
	}
	defer p.Close()

	return &RetryConsumer{
		consumer:     c,
		producer:     p,
		topic:        topic,
		maxRetries:   config.KafkaMaxRetry,
		delayMinutes: 30,
	}, nil
}

func (rc *RetryConsumer) Start(ctx context.Context) {
	err := rc.consumer.Subscribe(rc.topic, nil)
	if err != nil {
		log.Errorf("failed to subscribe to topic %s: %v", rc.topic, err)
		return
	}

	log.Println("Retry Consumer started...")

	for {
		message, err := rc.consumer.ReadMessage(-1)
		if err != nil {
			log.Errorf("Consumer error: %v\n", err)
			return
		}

		headers := parseHeaders(message.Headers)
		originalTimestamp, _ := strconv.ParseInt(headers["original_timestamp"], 10, 64)
		retryCount, _ := strconv.Atoi(headers["retry_count"])

		now := time.Now().Unix()
		if now-originalTimestamp < int64(rc.delayMinutes*60) {
			log.Printf("Skipping message, not enough delay yet: %v seconds remaining", int64(rc.delayMinutes*60)-(now-originalTimestamp))
			return
		}

		log.Printf("Processing retry #%d: %v", retryCount, message)

		orderService, err := service.NewOrder(message)
		if err != nil {
			log.Error("[ProcessMessage] - Error creating order service:", err)
			return
		}

		_, err = orderService.Process()
		if err != nil {
			if retryCount >= rc.maxRetries {
				log.Printf("Max retries reached. Sending to DLQ: %v", message)
				rc.producer.ProduceMessageDlq(ctx, message)
			} else {
				log.Printf("Reproducing to Retry Topic: %v", message)
				rc.producer.ProduceMessage(ctx, message, retryCount)
			}
		}

		log.Printf("Processed successfully: %v", message)
	}
}

func (rc *RetryConsumer) Close() {
	rc.consumer.Close()
}

func parseHeaders(headers []kafka.Header) map[string]string {
	result := make(map[string]string)
	for _, h := range headers {
		result[h.Key] = string(h.Value)
	}
	return result
}
