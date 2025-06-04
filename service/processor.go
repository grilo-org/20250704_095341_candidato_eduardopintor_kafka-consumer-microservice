package service

import (
	"bytes"
	"eduardopintor/kafka-consumer-microservice.git/config"
	"eduardopintor/kafka-consumer-microservice.git/dto"
	"encoding/json"
	"fmt"
	"net/http"

	kafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	log "github.com/sirupsen/logrus"
)

type Order struct {
	order *dto.Order
}

func NewOrder(message *kafka.Message) (*Order, error) {
	var order dto.Order
	err := json.Unmarshal(message.Value, &order)
	if err != nil {
		log.Error("[ProcessOrder] - Error unmarshalling message:", err)
		return nil, err
	}

	return &Order{
		order: &order,
	}, nil
}

func (o *Order) Process() (string, error) {
	// TODO: process the order

	if o.sendOrderToApi() != nil {
		return "", fmt.Errorf("failed to send order to API: %v", o.sendOrderToApi())
	}

	return "Order processed successfully", nil
}

func (o *Order) sendOrderToApi() error {
	orderBody, _ := json.Marshal(
		map[string]any{
			"ordemDeVenda": o.order.OrderID,
			"etapaAtual":   o.order.OrderStatus,
		},
	)

	req, err := http.NewRequest("PATCH", config.TargetServiceUrl, bytes.NewBuffer(orderBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	return nil
}
