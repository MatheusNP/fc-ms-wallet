package getaccountbalance

import (
	"testing"

	"github.com/MatheusNP/fc-ms-wallet/ms-account/internal/entity"
	"github.com/MatheusNP/fc-ms-wallet/ms-account/internal/usecase/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetAccountBalanceUseCase_Execute(t *testing.T) {
	account := entity.NewAccount(uuid.NewString(), 0.0)

	accountGateway := &mocks.AccountGatewayMock{}
	accountGateway.On("FindByID", account.ID).Return(account, nil)

	usecase := NewGetAccountBalanceUseCase(accountGateway)

	output, err := usecase.Execute(GetAccountBalanceInputDTO{
		AccountID: account.ID,
	})

	assert.Nil(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, output.Balance, account.Balance)

	accountGateway.AssertExpectations(t)
	accountGateway.AssertNumberOfCalls(t, "FindByID", 1)
}
