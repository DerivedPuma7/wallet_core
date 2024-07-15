package event

import (
	"github.com.br/derivedpuma7/wallet-core/pkg/events"
	"time"
)

type TransactionCreated struct {
	Name    string
  Datetime time.Time
	Payload interface{}
}

func NewTransactionCreated(datetime time.Time) *TransactionCreated {
	return &TransactionCreated{
		Name: "TransactionCreated",
    Datetime: datetime,
	}
}

func (t *TransactionCreated) GetName() string {
	return t.Name
}

func (t *TransactionCreated) GetDatetime() time.Time {
	return t.Datetime
}

func (t *TransactionCreated) GetPayload() interface{} {
	return t.Payload
}

func (t *TransactionCreated) SetPayload(payload interface{}) {
	t.Payload = payload
}

var _ events.EventInterface = (*TransactionCreated)(nil)
