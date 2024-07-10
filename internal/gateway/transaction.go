package gateway

import "github.com.br/derivedpuma7/wallet-core/internal/entity"

type TransactionGateway interface {
	Create(transaction *entity.Transaction)  error
}