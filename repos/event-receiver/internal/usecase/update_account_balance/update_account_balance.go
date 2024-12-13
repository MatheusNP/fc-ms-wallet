package updateaccountbalance

import (
	"github.com/MatheusNP/fc-ms-wallet/ms-account/internal/entity"
	"github.com/MatheusNP/fc-ms-wallet/ms-account/internal/gateway"
)

type UpdateAccountBalanceInputDTO struct {
	AccountID string  `json:"account_id"`
	Balance   float64 `json:"balance"`
}

type UpdateAccountBalanceUseCase struct {
	AccountGateway gateway.AccountGateway
}

func NewUpdateAccountBalanceUseCase(
	accountGateway gateway.AccountGateway,
) *UpdateAccountBalanceUseCase {
	return &UpdateAccountBalanceUseCase{
		AccountGateway: accountGateway,
	}
}

func (uc *UpdateAccountBalanceUseCase) Execute(input UpdateAccountBalanceInputDTO) error {
	var account *entity.Account

	account, err := uc.AccountGateway.FindByID(input.AccountID)
	if err != nil {
		// TODO: pegar erro exato
		if err.Error() != "sql: no rows in result set" {
			return err
		}

		account = entity.NewAccount(input.AccountID, input.Balance)

		err = uc.AccountGateway.Save(account)
		if err != nil {
			return err
		}
	}

	account.Balance = input.Balance

	return uc.AccountGateway.UpdateBalance(account)
}
