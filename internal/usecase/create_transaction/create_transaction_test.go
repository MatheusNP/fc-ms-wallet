package createtransaction

import (
	"testing"

	"github.com/MatheusNP/fc-ms-wallet/internal/entity"
	"github.com/MatheusNP/fc-ms-wallet/internal/event"
	"github.com/MatheusNP/fc-ms-wallet/pkg/events"
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

type mockTransactionGateway struct {
	mock.Mock
}

func (m *mockTransactionGateway) Create(transaction *entity.Transaction) error {
	args := m.Called(transaction)
	return args.Error(0)
}

func TestCreateTransaction(t *testing.T) {
	client1, _ := entity.NewClient("name", "email")
	account1 := entity.NewAccount(client1)
	client1.AddAccount(account1)
	account1.Credit(100.0)

	client2, _ := entity.NewClient("name", "email")
	account2 := entity.NewAccount(client2)
	client2.AddAccount(account2)
	account2.Credit(100.0)

	mockAccountGateway := &mockAccountGateway{}
	mockAccountGateway.On("FindByID", account1.ID).Return(account1, nil)
	mockAccountGateway.On("FindByID", account2.ID).Return(account2, nil)

	mockTransactionGateway := &mockTransactionGateway{}
	mockTransactionGateway.On("Create", mock.Anything).Return(nil)

	dispatcher := events.NewEventDispatcher()
	event := event.NewTransactionCreated()

	usecase := NewCreateTransactionUseCase(
		mockAccountGateway,
		mockTransactionGateway,
		dispatcher,
		event,
	)

	input := CreateTransactionInputDTO{
		AccountFromID: account1.ID,
		AccountToID:   account2.ID,
		Amount:        10.0,
	}

	output, err := usecase.Execute(input)

	assert.Nil(t, err)
	assert.NotNil(t, output)
	assert.NotEmpty(t, output.ID)

	mockAccountGateway.AssertExpectations(t)
	mockAccountGateway.AssertNumberOfCalls(t, "FindByID", 2)
	mockTransactionGateway.AssertExpectations(t)
	mockTransactionGateway.AssertNumberOfCalls(t, "Create", 1)
}
