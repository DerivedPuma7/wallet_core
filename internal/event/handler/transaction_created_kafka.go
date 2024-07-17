package handler

import (
	"fmt"
	"sync"

	"github.com.br/derivedpuma7/wallet-core/pkg/events"
	"github.com.br/derivedpuma7/wallet-core/pkg/kafka"
)

type TransactionCreatedKafkaHandler struct {
	Kafka *kafka.Producer
}

var _ events.EventHandlerInterface = (*TransactionCreatedKafkaHandler)(nil)

func NewTransactionCreatedKafkaHandler(kafka *kafka.Producer) *TransactionCreatedKafkaHandler {
	return &TransactionCreatedKafkaHandler{
		Kafka: kafka,
	}
}

func (h *TransactionCreatedKafkaHandler) Handle(event events.EventInterface, wg *sync.WaitGroup) {
	defer wg.Done()
  h.Kafka.Publish(event, nil, "transactions")
  fmt.Println(`TransactionCreatedKafkaHandler:`, event.GetPayload())
}
