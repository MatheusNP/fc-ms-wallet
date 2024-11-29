package createaccount

import (
	"testing"

	"github.com/MatheusNP/fc-ms-wallet/internal/entity"
	"github.com/MatheusNP/fc-ms-wallet/internal/usecase/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateAccountUseCase_Execute(t *testing.T) {
	client, _ := entity.NewClient("name", "email")

	clientGateway := &mocks.ClientGatewayMock{}
	clientGateway.On("FindByID", client.ID).Return(client, nil)

	accountGateway := &mocks.AccountGatewayMock{}
	accountGateway.On("Save", mock.Anything).Return(nil)

	usecase := NewCreateAccountUseCase(clientGateway, accountGateway)

	output, err := usecase.Execute(CreateAccountInputDTO{
		ClientID: client.ID,
	})

	assert.Nil(t, err)
	assert.NotNil(t, output)
	assert.NotEmpty(t, output.ID)

	clientGateway.AssertExpectations(t)
	clientGateway.AssertNumberOfCalls(t, "FindByID", 1)
	accountGateway.AssertExpectations(t)
	accountGateway.AssertNumberOfCalls(t, "Save", 1)
}
