package mocks

import (
	"github.com.br/derivedpuma7/wallet-core/internal/entity"
	"github.com.br/derivedpuma7/wallet-core/internal/gateway"
	"github.com/stretchr/testify/mock"
)

type TransactionGatewayMock struct {
	mock.Mock
}

var _ gateway.TransactionGateway = (*TransactionGatewayMock)(nil)

func (m *TransactionGatewayMock) Create(transaction *entity.Transaction) error {
	args := m.Called(transaction)
	return args.Error(0)
}
