package createclient

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

func TestCreateClientUseCase_Execute(t *testing.T) {
	m := &ClientGatewayMock{}
	m.On("Save", mock.Anything).Return(nil)
	uc := NewCreateClientUseCase(m)

	output, err := uc.Execute(CreateClientInputDto{
		Name: "any name",
		Email: "any email",
	})

	assert.Nil(t, err)
	assert.NotNil(t, output)
	assert.NotEmpty(t, output.ID)
	assert.Equal(t, "any name", output.Name)
	assert.Equal(t, "any email", output.Email)
	m.AssertExpectations(t)
	m.AssertNumberOfCalls(t, "Save", 1)
}