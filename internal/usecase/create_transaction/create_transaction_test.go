package createtransaction

import (
	"context"
	"testing"
	"time"

	"github.com.br/derivedpuma7/wallet-core/internal/entity"
	"github.com.br/derivedpuma7/wallet-core/internal/event"
	"github.com.br/derivedpuma7/wallet-core/internal/usecase/mocks"
	"github.com.br/derivedpuma7/wallet-core/pkg/events"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func createUowMock() *mocks.UowMock {
  mockUow := &mocks.UowMock{}
  mockUow.On("Do", mock.Anything, mock.Anything).Return(nil)
  return mockUow
}

func TestCreateTransactionUseCase_Execute(t *testing.T) {
	client1, _ := entity.NewClient("any name 1", "any email 1")
	account1, _ := entity.NewAccount(client1)
	account1.Credit(1000)
	client2, _ := entity.NewClient("any name 2", "any email 2")
	account2, _ := entity.NewAccount(client2)
	account2.Credit(1000)
  mockUow := createUowMock()
  ctx := context.Background()
	dispatcher := events.NewEventDispatcher()
	eventTransaction := event.NewTransactionCreated(time.Now())
  eventBalance := event.NewBalanceUpdated()

	uc := NewCreateTransactionUseCase(mockUow, dispatcher, eventTransaction, eventBalance)

	output, err := uc.Execute(
    ctx,
		CreateTransactionInputDto{
      AccountIDFrom: account1.ID,
      AccountIDTo: account2.ID,
      Amount: 100,
    },
  )

	assert.Nil(t, err)
	assert.NotNil(t, output)
	mockUow.AssertExpectations(t)
	mockUow.AssertNumberOfCalls(t, "Do", 1)
}
