package main

import (
	"context"
	"eduardopintor/kafka-consumer-microservice.git/config"
	"eduardopintor/kafka-consumer-microservice.git/consumer"

	log "github.com/sirupsen/logrus"
)

func main() {
	ctx := context.Background()
	if err := config.Init(); err != nil {
		log.Error("Error initializing configuration: ", err)
		return
	}

	msgConsumer, err := consumer.NewConsumer(config.KafkaTopic)
	if err != nil {
		log.Error("Error creating consumer: ", err)
		return
	}
	defer msgConsumer.Close()

	log.Info("Consumer created successfully, starting to consume messages...")

	go msgConsumer.Start(ctx)

	retryConsumer, err := consumer.NewRetryConsumer(config.KafkaTopicRetry)
	if err != nil {
		log.Error("Error creating retry consumer: ", err)
		return
	}
	defer retryConsumer.Close()

	log.Info("Retry Consumer created successfully, starting to consume retry messages...")
	// go retryConsumer.Start(ctx)

	<-ctx.Done()
}
