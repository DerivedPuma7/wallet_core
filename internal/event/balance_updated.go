package event

import (
	"github.com.br/derivedpuma7/wallet-core/pkg/events"
	"time"
)

type BalanceUpdated struct {
  Name    string
	Payload interface{}
}

func NewBalanceUpdated() *BalanceUpdated {
	return &BalanceUpdated{
		Name: "BalanceUpdated",
	}
}

var _ events.EventInterface = (*BalanceUpdated)(nil)

func (b *BalanceUpdated) GetDatetime() time.Time {
	return time.Now()
}

func (b *BalanceUpdated) GetName() string {
	return b.Name
}

func (b *BalanceUpdated) GetPayload() interface{} {
	return b.Payload
}

func (b *BalanceUpdated) SetPayload(payload interface{}) {
	b.Payload = payload
}
