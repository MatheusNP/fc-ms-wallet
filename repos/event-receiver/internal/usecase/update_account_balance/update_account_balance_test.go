package updateaccountbalance

import (
	"testing"

	"github.com/MatheusNP/fc-ms-wallet/ms-account/internal/entity"
	"github.com/MatheusNP/fc-ms-wallet/ms-account/internal/usecase/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUpdateAccountBalanceUseCase_Execute(t *testing.T) {
	account := entity.NewAccount(uuid.NewString(), 0.0)

	accountGateway := &mocks.AccountGatewayMock{}
	accountGateway.On("FindByID", account.ID).Return(account, nil)
	accountGateway.On("UpdateBalance", account).Return(nil)

	usecase := NewUpdateAccountBalanceUseCase(accountGateway)

	err := usecase.Execute(UpdateAccountBalanceInputDTO{
		AccountID: account.ID,
		Balance:   account.Balance,
	})

	assert.Nil(t, err)

	accountGateway.AssertExpectations(t)
	accountGateway.AssertNumberOfCalls(t, "FindByID", 1)
	accountGateway.AssertNumberOfCalls(t, "UpdateBalance", 1)
}
