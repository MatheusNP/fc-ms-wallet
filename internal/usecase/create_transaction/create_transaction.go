package createtransaction

import (
	"context"

	"github.com/MatheusNP/fc-ms-wallet/internal/entity"
	"github.com/MatheusNP/fc-ms-wallet/internal/gateway"
	"github.com/MatheusNP/fc-ms-wallet/pkg/events"
	"github.com/MatheusNP/fc-ms-wallet/pkg/uow"
)

type CreateTransactionInputDTO struct {
	AccountFromID string  `json:"account_from_id"`
	AccountToID   string  `json:"account_to_id"`
	Amount        float64 `json:"amount"`
}

type CreateTransactionOutputDTO struct {
	ID            string  `json:"id"`
	AccountFromID string  `json:"account_from_id"`
	AccountToID   string  `json:"account_to_id"`
	Amount        float64 `json:"amount"`
}

type BalanceUpdatedOutputDTO struct {
	AccountFromID      string  `json:"account_from_id"`
	AccountToID        string  `json:"account_to_id"`
	BalanceAccountFrom float64 `json:"balance_account_from"`
	BalanceAccountTo   float64 `json:"balance_account_to"`
}

type CreateTransactionUseCase struct {
	Uow                uow.UowInterface
	EventDispatcher    events.EventDispatcherInterface
	TransactionCreated events.EventInterface
	BalanceUpdated     events.EventInterface
}

func NewCreateTransactionUseCase(
	uow uow.UowInterface,
	eventDispatcher events.EventDispatcherInterface,
	transactionCreated events.EventInterface,
	balanceUpdated events.EventInterface,
) *CreateTransactionUseCase {
	return &CreateTransactionUseCase{
		Uow:                uow,
		EventDispatcher:    eventDispatcher,
		TransactionCreated: transactionCreated,
		BalanceUpdated:     balanceUpdated,
	}
}

func (uc *CreateTransactionUseCase) Execute(ctx context.Context, input CreateTransactionInputDTO) (*CreateTransactionOutputDTO, error) {

	output := &CreateTransactionOutputDTO{}

	balanceUpdatedOutput := &BalanceUpdatedOutputDTO{}

	err := uc.Uow.Do(ctx, func(_ *uow.Uow) error {
		accountRepository := uc.getAccountRepository(ctx)
		transactionRepository := uc.getTransactionRepository(ctx)

		accountFrom, err := accountRepository.FindByID(input.AccountFromID)
		if err != nil {
			return err
		}

		accountTo, err := accountRepository.FindByID(input.AccountToID)
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
		output.AccountFromID = transaction.AccountFrom.ID
		output.AccountToID = transaction.AccountTo.ID
		output.Amount = transaction.Amount

		balanceUpdatedOutput.AccountFromID = accountFrom.ID
		balanceUpdatedOutput.AccountToID = accountTo.ID
		balanceUpdatedOutput.BalanceAccountFrom = accountFrom.Balance
		balanceUpdatedOutput.BalanceAccountTo = accountTo.Balance

		return nil
	})
	if err != nil {
		return nil, err
	}

	uc.TransactionCreated.SetPayload(output)
	uc.EventDispatcher.Dispatch(uc.TransactionCreated)

	uc.BalanceUpdated.SetPayload(balanceUpdatedOutput)
	uc.EventDispatcher.Dispatch(uc.BalanceUpdated)

	return output, nil
}

func (uc *CreateTransactionUseCase) getAccountRepository(ctx context.Context) gateway.AccountGateway {
	repo, err := uc.Uow.GetRepository(ctx, "AccountDB")
	if err != nil {
		panic(err)
	}

	return repo.(gateway.AccountGateway)
}

func (uc *CreateTransactionUseCase) getTransactionRepository(ctx context.Context) gateway.TransactionGateway {
	repo, err := uc.Uow.GetRepository(ctx, "TransactionDB")
	if err != nil {
		panic(err)
	}

	return repo.(gateway.TransactionGateway)
}
