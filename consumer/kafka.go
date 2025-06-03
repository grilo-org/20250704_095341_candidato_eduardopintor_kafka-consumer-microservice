package consumer

import (
	"eduardopintor/kafka-consumer-microservice.git/config"
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	log "github.com/sirupsen/logrus"
)

// ConsumeMessages reads messages from a Kafka topic and processes them
func ConsumeMessages() error {
	// Create a new Kafka consumer
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": config.KafkaBootstrapServers,
		"group.id":          config.KafkaGroupId,
		"auto.offset.reset": "earliest",
		"sasl.username":     config.KafkaUsername,
		"sasl.password":     config.KafkaPassword,
		"sasl.mechanism":    config.KafkaSaslMechanism,
		"security.protocol": config.KafkaSecurityProtocol,
	})

	if err != nil {
		return fmt.Errorf("failed to create Kafka consumer: %v", err)
	}

	defer consumer.Close()

	err = consumer.Subscribe(config.KafkaTopic, nil)
	if err != nil {
		return fmt.Errorf("failed to subscribe to topic %s: %v", config.KafkaTopic, err)
	}

	for {
		msg, err := consumer.ReadMessage(-1)
		if err != nil {
			log.Errorf("Error reading message: %v", err)
			continue
		}

		log.Infof("Received message: %s", string(msg.Value))

		time.Sleep(1 * time.Second) // Simulate processing time
	}
}
