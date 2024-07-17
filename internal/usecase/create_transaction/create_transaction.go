package createtransaction

import (
	"context"

	"github.com.br/derivedpuma7/wallet-core/internal/entity"
	"github.com.br/derivedpuma7/wallet-core/internal/gateway"
	"github.com.br/derivedpuma7/wallet-core/pkg/events"
	"github.com.br/derivedpuma7/wallet-core/pkg/uow"
)

type CreateTransactionInputDto struct {
	AccountIDFrom string
	AccountIDTo string
	Amount float64
}

type CreateClientOutputDto struct {
	ID string
  AccountIdFrom string
  AccountIdTo string
  Amount float64
}

type CreateTransactionUseCase struct {
  Uow uow.UowInterface
	EventDispatcher events.EventDispatcherInterface
	TransactionCreated events.EventInterface
}

func NewCreateTransactionUseCase(
  uow uow.UowInterface,
	eventDispatcher events.EventDispatcherInterface,
	transactionCreated events.EventInterface,
) *CreateTransactionUseCase {
	return &CreateTransactionUseCase{
    Uow: uow,
		EventDispatcher: eventDispatcher,
		TransactionCreated: transactionCreated,
	}
}

func (uc *CreateTransactionUseCase) getAccountRepository(ctx context.Context) gateway.AccountGateway {
  repo, err := uc.Uow.GetRepository(ctx, "AccountDb")
  if err != nil {
    panic(err)
  }
  return repo.(gateway.AccountGateway)
}

func (uc *CreateTransactionUseCase) getTransactionRepository(ctx context.Context) gateway.TransactionGateway {
  repo, err := uc.Uow.GetRepository(ctx, "TransactionDb")
  if err != nil {
    panic(err)
  }
  return repo.(gateway.TransactionGateway)
}

func (uc *CreateTransactionUseCase) Execute(ctx context.Context, input CreateTransactionInputDto) (*CreateClientOutputDto, error) {
  output := &CreateClientOutputDto{}
  err := uc.Uow.Do(ctx, func(_ *uow.Uow) error {
    accountRepository := uc.getAccountRepository(ctx)
    transactionRepository := uc.getTransactionRepository(ctx)
    
    accountFrom, err := accountRepository.FindById(input.AccountIDFrom)
    if err != nil {
      return err
    }
    accountTo, err := accountRepository.FindById(input.AccountIDTo)
    if err != nil {
      return err
    }
    transaction, err := entity.NewTransaction(accountFrom, accountTo, input.Amount)
    if err != nil {
      return err
    }
    err = accountRepository.UpdateBalance(accountFrom)
    if err != nil {
      return err
    }
    err = accountRepository.UpdateBalance(accountTo)
    if err != nil {
      return err
    }
    err = transactionRepository.Create(transaction)
    if err != nil {
      return err
    }
    output.ID = transaction.ID
    output.AccountIdFrom = accountFrom.ID
    output.AccountIdTo = accountTo.ID
    output.Amount = transaction.Amount
    return nil
  })
  if err != nil {
    return nil, err
  }
	uc.TransactionCreated.SetPayload(output)
	uc.EventDispatcher.Dispatch(uc.TransactionCreated)
	return output, nil
}
