package createaccount

import (
	"testing"

	"github.com/MatheusNP/fc-ms-wallet/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockAccountGateway struct {
	mock.Mock
}

func (m *mockAccountGateway) Save(account *entity.Account) error {
	args := m.Called(account)
	return args.Error(0)
}

func (m *mockAccountGateway) FindByID(id string) (*entity.Account, error) {
	args := m.Called(id)
	return args.Get(0).(*entity.Account), args.Error(1)
}

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

func TestCreateAccountUseCase_Execute(t *testing.T) {
	client, _ := entity.NewClient("name", "email")

	clientGateway := &mockClientGateway{}
	clientGateway.On("Get", client.ID).Return(client, nil)

	accountGateway := &mockAccountGateway{}
	accountGateway.On("Save", mock.Anything).Return(nil)

	usecase := NewCreateAccountUseCase(clientGateway, accountGateway)

	output, err := usecase.Execute(CreateAccountInputDTO{
		ClientID: client.ID,
	})

	assert.Nil(t, err)
	assert.NotNil(t, output)
	assert.NotEmpty(t, output.ID)

	clientGateway.AssertExpectations(t)
	clientGateway.AssertNumberOfCalls(t, "Get", 1)
	accountGateway.AssertExpectations(t)
	accountGateway.AssertNumberOfCalls(t, "Save", 1)
}
