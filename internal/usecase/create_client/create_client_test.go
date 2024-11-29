package createclient

import (
	"testing"

	"github.com/MatheusNP/fc-ms-wallet/internal/usecase/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateClient(t *testing.T) {
	clientGateway := &mocks.ClientGatewayMock{}
	clientGateway.On("Save", mock.Anything).Return(nil)

	usecase := NewCreateClientUseCase(clientGateway)
	output, err := usecase.Execute(CreateClientInputDTO{
		Name:  "name",
		Email: "email",
	})

	assert.Nil(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, "name", output.Name)
	assert.Equal(t, "email", output.Email)

	clientGateway.AssertExpectations(t)
	clientGateway.AssertNumberOfCalls(t, "Save", 1)
}
