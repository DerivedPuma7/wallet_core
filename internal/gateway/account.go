package gateway

import "github.com.br/derivedpuma7/wallet-core/internal/entity"

type AccountGateway interface {
	FindById(id string) (*entity.Account, error)
	Save(account *entity.Account) error
  UpdateBalance(account *entity.Account) error
}
