package handler

import (
	"fmt"
	"sync"

	"github.com.br/derivedpuma7/wallet-core/pkg/events"
	"github.com.br/derivedpuma7/wallet-core/pkg/kafka"
)

type BalanceUpdatedKafkaHandler struct {
	Kafka *kafka.Producer
}

var _ events.EventHandlerInterface = (*BalanceUpdatedKafkaHandler)(nil)

func NewBalanceUpdatedKafkaHandler(kafka *kafka.Producer) *BalanceUpdatedKafkaHandler {
	return &BalanceUpdatedKafkaHandler{
		Kafka: kafka,
	}
}

func (b *BalanceUpdatedKafkaHandler) Handle(event events.EventInterface, wg *sync.WaitGroup) {
	defer wg.Done()
  b.Kafka.Publish(event, nil, "balances")
  fmt.Println(`BalanceUpdatedKafkaHandler:`, event.GetPayload())
}
