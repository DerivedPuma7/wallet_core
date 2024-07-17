package createaccount

import (
	"testing"

	"github.com.br/derivedpuma7/wallet-core/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type ClientGatewayMock struct {
	mock.Mock
}

func (mock *ClientGatewayMock) Get(id string) (*entity.Client, error) {
	args := mock.Called(id)
	return args.Get(0).(*entity.Client), args.Error(1)
}

func (mock *ClientGatewayMock) Save(client *entity.Client) error {
	args := mock.Called(client)
	return args.Error(0)
}

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

func (mock *AccountGatewayMock) UpdateBalance(account *entity.Account) error {
	args := mock.Called(account)
	return args.Error(0)
}

func TestCreateAccountUseCase_Execute(t *testing.T) {
	client, _ := entity.NewClient("any name", "any email")
	clientMock := &ClientGatewayMock {}
	clientMock.On("Get", client.ID).Return(client, nil)
	accountMock := &AccountGatewayMock{}
	accountMock.On("Save", mock.Anything).Return(nil)
	
	uc := NewCreateAccountUseCase(accountMock, clientMock)
	output, err := uc.Execute(CreateAccountInputDto{
		ClientId: client.ID,
	})

	assert.Nil(t, err)
	assert.NotNil(t, output)
	assert.NotEmpty(t, output.ID)
	clientMock.AssertExpectations(t)
	clientMock.AssertNumberOfCalls(t, "Get", 1)
	accountMock.AssertExpectations(t)
	accountMock.AssertNumberOfCalls(t, "Save", 1)
}
