package createtransaction
import (
	"testing"
	"time"

	"github.com.br/derivedpuma7/wallet-core/internal/entity"
	"github.com.br/derivedpuma7/wallet-core/internal/event"
	"github.com.br/derivedpuma7/wallet-core/pkg/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type AccountGatewayMock struct {
	mock.Mock
}

func (mock *AccountGatewayMock) FindById(id string) (*entity.Account, error) {
	args := mock.Called(id)
	return args.Get(0).(*entity.Account), args.Error(1)
}

func (mock *AccountGatewayMock) Save(account *entity.Account) error {
	args := mock.Called(account)
	return args.Error(0)
}

type TransactionGatewayMock struct {
	mock.Mock
}

func (mock *TransactionGatewayMock) Create(transaction *entity.Transaction) error {
	args := mock.Called(transaction)
	return args.Error(0)
}

func createGatewaysMock(account1 *entity.Account, account2 *entity.Account) (*AccountGatewayMock, *TransactionGatewayMock) {
	mockAccount := &AccountGatewayMock{}
	mockAccount.On("FindById", account1.ID).Return(account1, nil)
	mockAccount.On("FindById", account2.ID).Return(account2, nil)
	mockTransaction := &TransactionGatewayMock{}
	mockTransaction.On("Create", mock.Anything).Return(nil)
	return mockAccount, mockTransaction
}

func TestCreateTransactionUseCase_Execute(t *testing.T) {
	client1, _ := entity.NewClient("any name 1", "any email 1")
	account1, _ := entity.NewAccount(client1)
	account1.Credit(1000)
	client2, _ := entity.NewClient("any name 2", "any email 2")
	account2, _ := entity.NewAccount(client2)
	account2.Credit(1000)
	mockAccount, mockTransaction := createGatewaysMock(account1, account2)
	dispatcher := events.NewEventDispatcher()
	event := event.NewTransactionCreated(time.Now())

	uc := NewCreateTransactionUseCase(mockTransaction, mockAccount, dispatcher, event)

	output, err := uc.Execute(
		CreateTransactionInputDto{
		AccountIDFrom: account1.ID,
		AccountIDTo: account2.ID,
		Amount: 100,
	})

	assert.Nil(t, err)
	assert.NotNil(t, output)
	mockAccount.AssertExpectations(t)
	mockAccount.AssertNumberOfCalls(t, "FindById", 2)
	mockTransaction.AssertExpectations(t)
	mockTransaction.AssertNumberOfCalls(t, "Create", 1)
}
