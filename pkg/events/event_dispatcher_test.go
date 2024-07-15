package events

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestSuite(t *testing.T) {
  suite.Run(t, new(EventDispatcherTestSuite))
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Register() {
  eventName := suite.event.GetName()

  err := suite.eventDispatcher.Register(eventName, &suite.handler)
  suite.Nil(err)
  suite.Equal(1, len(suite.eventDispatcher.handlers[eventName]))
  assert.Equal(suite.T(), &suite.handler, suite.eventDispatcher.handlers[eventName][0])

  err = suite.eventDispatcher.Register(eventName, &suite.handler2)
  suite.Nil(err)
  suite.Equal(2, len(suite.eventDispatcher.handlers[eventName]))
  assert.Equal(suite.T(), &suite.handler2, suite.eventDispatcher.handlers[eventName][1])
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Register_WithSameHandler() {
  eventName := suite.event.GetName()

  suite.eventDispatcher.Register(eventName, &suite.handler)
  err := suite.eventDispatcher.Register(eventName, &suite.handler)

  suite.NotNil(err)
  suite.Equal(ErrHandlerAlreadyRegistered, err)
  suite.Equal(1, len(suite.eventDispatcher.handlers[eventName]))
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Clear() {
  eventName := suite.event.GetName()
  eventName2 := suite.event2.GetName()
  
  // Event 1
  err := suite.eventDispatcher.Register(eventName, &suite.handler)
  suite.Nil(err)
  suite.Equal(1, len(suite.eventDispatcher.handlers[eventName]))
  err = suite.eventDispatcher.Register(eventName, &suite.handler2)
  suite.Nil(err)
  suite.Equal(2, len(suite.eventDispatcher.handlers[eventName]))

  // Event 2
  err = suite.eventDispatcher.Register(eventName2, &suite.handler3)
  suite.Nil(err)
  suite.Equal(1, len(suite.eventDispatcher.handlers[eventName2]))

  suite.eventDispatcher.Clear()

  suite.Equal(0, len(suite.eventDispatcher.handlers))
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Has() {
  eventName := suite.event.GetName()
  suite.eventDispatcher.Register(eventName, &suite.handler)
  suite.eventDispatcher.Register(eventName, &suite.handler2)

  has1 := suite.eventDispatcher.Has(eventName, &suite.handler)
  has2 := suite.eventDispatcher.Has(eventName, &suite.handler2)
  hasnot := suite.eventDispatcher.Has(eventName, &suite.handler3)

  suite.Equal(true, has1)
  suite.Equal(true, has2)
  suite.Equal(false, hasnot)
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Remove() {
  eventName := suite.event.GetName()
  suite.eventDispatcher.Register(eventName, &suite.handler)
  suite.eventDispatcher.Register(eventName, &suite.handler2)

  suite.Equal(2, len(suite.eventDispatcher.handlers[eventName]))

  suite.eventDispatcher.Remove(eventName, &suite.handler)

  suite.Equal(1, len(suite.eventDispatcher.handlers[eventName]))
  suite.Equal(&suite.handler2, suite.eventDispatcher.handlers[eventName][0])
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Dispatch() {
  eHandler := &MockHandler{}
  eHandler.On("Handle", &suite.event)
  suite.eventDispatcher.Register(suite.event.GetName(), eHandler)

  suite.eventDispatcher.Dispatch(&suite.event)

  eHandler.AssertExpectations(suite.T())
  eHandler.AssertNumberOfCalls(suite.T(), "Handle", 1)
}
