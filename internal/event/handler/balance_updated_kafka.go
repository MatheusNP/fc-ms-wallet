package handler

import (
	"fmt"
	"sync"

	"github.com/MatheusNP/fc-ms-wallet/pkg/events"
	"github.com/MatheusNP/fc-ms-wallet/pkg/kafka"
)

type UpdateBalanceKafkaHandler struct {
	Kafka *kafka.Producer
}

func NewUpdateBalanceKafkaHandler(kafka *kafka.Producer) *UpdateBalanceKafkaHandler {
	return &UpdateBalanceKafkaHandler{Kafka: kafka}
}

func (h *UpdateBalanceKafkaHandler) Handle(msg events.EventInterface, wg *sync.WaitGroup) {
	defer wg.Done()
	h.Kafka.Publish(msg, nil, "balances")
	fmt.Println("UpdateBalanceKafkaHandler: ", msg.GetPayload())
}
