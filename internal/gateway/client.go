package gateway

import "github.com.br/derivedpuma7/wallet-core/internal/entity"

type ClientGateway interface {
	Get(id string) (*entity.Client, error)
	Save(client *entity.Client) error
}