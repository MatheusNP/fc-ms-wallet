package createaccount

import (
	"fmt"

	"github.com/MatheusNP/fc-ms-wallet/internal/entity"
	"github.com/MatheusNP/fc-ms-wallet/internal/gateway"
)

type CreateAccountInputDTO struct {
	ClientID string `json:"client_id"`
}

type CreateAccountOutputDTO struct {
	ID string `json:"id"`
}

type CreateAccountUseCase struct {
	ClientGateway  gateway.ClientGateway
	AccountGateway gateway.AccountGateway
}

func NewCreateAccountUseCase(
	clientGateway gateway.ClientGateway,
	accountGateway gateway.AccountGateway,
) *CreateAccountUseCase {
	return &CreateAccountUseCase{
		ClientGateway:  clientGateway,
		AccountGateway: accountGateway,
	}
}

func (uc *CreateAccountUseCase) Execute(input CreateAccountInputDTO) (*CreateAccountOutputDTO, error) {
	client, err := uc.ClientGateway.FindByID(input.ClientID)
	if err != nil {
		fmt.Println(input.ClientID)
		return nil, err
	}

	account := entity.NewAccount(client)

	err = uc.AccountGateway.Save(account)
	if err != nil {
		return nil, err
	}

	return &CreateAccountOutputDTO{
		ID: account.ID,
	}, nil
}
