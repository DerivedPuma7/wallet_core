package createaccount

import (
	"github.com.br/derivedpuma7/wallet-core/internal/entity"
	"github.com.br/derivedpuma7/wallet-core/internal/gateway"
)

type CreateAccountInputDto struct {
	ClientId string
}

type CreateAccountOutputDto struct {
	ID string
}

type CreateAccountUseCase struct {
	AccountGateway gateway.AccountGateway
	ClientGateway gateway.ClientGateway
}

func NewCreateAccountUseCase(accountGateway gateway.AccountGateway, clientGateway gateway.ClientGateway) *CreateAccountUseCase {
	return &CreateAccountUseCase{
		AccountGateway: accountGateway,
		ClientGateway: clientGateway,
	}
}

func(uc *CreateAccountUseCase) Execute(input CreateAccountInputDto) (*CreateAccountOutputDto, error) {
	client, err := uc.ClientGateway.Get(input.ClientId)
	if err != nil {
		return nil, err
	}
	account, _ := entity.NewAccount(client)
	err = uc.AccountGateway.Save(account)
	if err != nil {
		return nil, err
	}
	return &CreateAccountOutputDto{
		ID: account.ID,
	}, nil
}