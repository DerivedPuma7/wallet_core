package mocks

import (
	"github.com.br/derivedpuma7/wallet-core/internal/entity"
	"github.com.br/derivedpuma7/wallet-core/internal/gateway"
	"github.com/stretchr/testify/mock"
)

type ClientGatewayMock struct {
	mock.Mock
}

var _ gateway.ClientGateway = (*ClientGatewayMock)(nil)

func (m *ClientGatewayMock) Save(client *entity.Client) error {
	args := m.Called(client)
	return args.Error(0)
}

func (m *ClientGatewayMock) Get(id string) (*entity.Client, error) {
	args := m.Called(id)
	return args.Get(0).(*entity.Client), args.Error(1)
}
