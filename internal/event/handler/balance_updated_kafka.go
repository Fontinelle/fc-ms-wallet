package handler

import (
	"fmt"
	"sync"

	"github.com/fontinelle/fc-ms-wallet/pkg/events"
	"github.com/fontinelle/fc-ms-wallet/pkg/kafka"
)

type UpdatedBalanceKafkaHandler struct {
	Kafka *kafka.Producer
}

func NewUpdateBalanceKafkaHandler(kafka *kafka.Producer) *UpdatedBalanceKafkaHandler {
	return &UpdatedBalanceKafkaHandler{
		Kafka: kafka,
	}
}

func (h *UpdatedBalanceKafkaHandler) Handle(message events.EventInterface, wg *sync.WaitGroup) {
	defer wg.Done()
	h.Kafka.Publish(message, nil, "balances")
	fmt.Println("UpdatedBalanceKafkaHandler called")
}
