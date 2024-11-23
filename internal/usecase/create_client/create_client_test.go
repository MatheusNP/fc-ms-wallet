package createclient

import (
	"testing"

	"github.com/MatheusNP/fc-ms-wallet/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockClientGateway struct {
	mock.Mock
}

func (m *mockClientGateway) Save(client *entity.Client) error {
	args := m.Called(client)
	return args.Error(0)
}

func (m *mockClientGateway) Get(id string) (*entity.Client, error) {
	args := m.Called(id)
	return args.Get(0).(*entity.Client), args.Error(1)
}

func TestCreateClient(t *testing.T) {
	clientGateway := &mockClientGateway{}
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
