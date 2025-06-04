package producer

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"eduardopintor/kafka-consumer-microservice.git/config"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	log "github.com/sirupsen/logrus"
)

type Producer struct {
	producer *kafka.Producer
	topic    string
}

func NewProducer(topic string) (*Producer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": config.KafkaBootstrapServers,
		"security.protocol": config.KafkaSecurityProtocol,
		"sasl.mechanisms":   config.KafkaSaslMechanism,
		"sasl.username":     config.KafkaUsername,
		"sasl.password":     config.KafkaPassword,
	})
	if err != nil {
		return nil, err
	}

	return &Producer{
		producer: p,
		topic:    topic, // <-- Add this line
	}, nil
}

func (rp *Producer) ProduceMessage(ctx context.Context, payload *kafka.Message, currentRetry int) error {
	// Marshal only the payload value, not the whole message
	// valueBytes, err := json.Marshal(payload.Value)
	// if err != nil {
	// 	return err
	// }

	originalTimestamp := strconv.FormatInt(time.Now().Unix(), 10)
	retryCountStr := strconv.Itoa(currentRetry + 1)

	message := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &rp.topic,
			Partition: kafka.PartitionAny, // Ensure partition is set
		},
		Value: payload.Value,
		Headers: []kafka.Header{
			{Key: "retry_count", Value: []byte(retryCountStr)},
			{Key: "original_timestamp", Value: []byte(originalTimestamp)},
		},
	}

	log.Printf("Producing message to topic: %s retry_count=%s, timestamp=%s", rp.topic, retryCountStr, originalTimestamp)

	err := rp.producer.Produce(message, nil)
	if err != nil {
		log.Errorf("Failed to produce message: %v", err)
		return err
	}

	// Wait for delivery report for this message
	e := <-rp.producer.Events()
	m, ok := e.(*kafka.Message)
	if !ok {
		log.Errorf("Failed to cast event to *kafka.Message")
		return fmt.Errorf("failed to cast event to *kafka.Message")
	}

	if m.TopicPartition.Error != nil {
		log.Printf("Failed to deliver message: %v", m.TopicPartition.Error)
		return m.TopicPartition.Error
	}

	log.Printf("Message delivered to %v", m.TopicPartition)
	return nil
}

func (rp *Producer) ProduceMessageDlq(ctx context.Context, payload *kafka.Message) error {
	rp.topic = config.KafkaTopicDlq

	return rp.ProduceMessage(ctx, payload, 0)
}

func (rp *Producer) Close() {
	rp.producer.Close()
}
