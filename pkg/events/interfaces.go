package events

import "time"

type EventInterface interface {
	GetName() string
	GetDatetime() time.Time
	GetPayload() interface{}
  SetPayload(payload interface{})
}

type EventHandlerInterface interface {
  Handle(event EventInterface)
}

type EventDispatcherInterface interface {
  Register(eventName string, handler EventHandlerInterface) error
  Remove(eventName string, handler EventHandlerInterface) error
  Dispatch(event EventInterface) error
  Has(eventName string, handler EventHandlerInterface) bool
  Clear() error
}