package getaccountbalance

import "github.com/MatheusNP/fc-ms-wallet/ms-account/internal/gateway"

type GetAccountBalanceInputDTO struct {
	AccountID string `json:"account_id"`
}

type GetAccountBalanceOutputDTO struct {
	Balance float64 `json:"balance"`
}

type GetAccountBalanceUseCase struct {
	AccountGateway gateway.AccountGateway
}

func NewGetAccountBalanceUseCase(
	accountGateway gateway.AccountGateway,
) *GetAccountBalanceUseCase {
	return &GetAccountBalanceUseCase{
		AccountGateway: accountGateway,
	}
}

func (uc *GetAccountBalanceUseCase) Execute(input GetAccountBalanceInputDTO) (*GetAccountBalanceOutputDTO, error) {
	account, err := uc.AccountGateway.FindByID(input.AccountID)
	if err != nil {
		return nil, err
	}

	return &GetAccountBalanceOutputDTO{
		Balance: account.Balance,
	}, nil
}
