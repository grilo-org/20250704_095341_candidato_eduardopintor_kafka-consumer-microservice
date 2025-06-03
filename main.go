package main

import (
	"eduardopintor/kafka-consumer-microservice.git/config"
	"eduardopintor/kafka-consumer-microservice.git/consumer"

	log "github.com/sirupsen/logrus"
)

func main() {
	if err := config.Init(); err != nil {
		log.Error("Error initializing configuration: ", err)
		return
	}

	if err := consumer.ConsumeMessages(); err != nil {
		log.Error("Error consuming messages: ", err)
		return
	}

}
