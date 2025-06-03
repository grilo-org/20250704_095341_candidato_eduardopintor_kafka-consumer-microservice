package main

import (
	"eduardopintor/kafka-consumer-microservice.git/config"
	"fmt"

	log "github.com/sirupsen/logrus"
)

func main() {
	fmt.Println("Consumer Up!")

	if err := config.Init(); err != nil {
		log.Error("Error initializing configuration: ", err)
		return
	}

}
