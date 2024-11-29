package createtransaction

import (
	"context"
	"testing"

	"github.com/MatheusNP/fc-ms-wallet/internal/entity"
	"github.com/MatheusNP/fc-ms-wallet/internal/event"
	"github.com/MatheusNP/fc-ms-wallet/internal/usecase/mocks"
	"github.com/MatheusNP/fc-ms-wallet/pkg/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateTransaction(t *testing.T) {
	client1, _ := entity.NewClient("name", "email")
	account1 := entity.NewAccount(client1)
	client1.AddAccount(account1)
	account1.Credit(100.0)

	client2, _ := entity.NewClient("name", "email")
	account2 := entity.NewAccount(client2)
	client2.AddAccount(account2)
	account2.Credit(100.0)

	mockUow := &mocks.UowMock{}
	mockUow.On("Do", mock.Anything, mock.Anything).Return(nil)

	dispatcher := events.NewEventDispatcher()
	event := event.NewTransactionCreated()
	ctx := context.Background()

	usecase := NewCreateTransactionUseCase(
		mockUow,
		dispatcher,
		event,
	)

	input := CreateTransactionInputDTO{
		AccountFromID: account1.ID,
		AccountToID:   account2.ID,
		Amount:        10.0,
	}

	output, err := usecase.Execute(ctx, input)

	assert.Nil(t, err)
	assert.NotNil(t, output)

	mockUow.AssertExpectations(t)
	mockUow.AssertNumberOfCalls(t, "Do", 1)
}
