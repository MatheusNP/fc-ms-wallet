package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateNewTransaction(t *testing.T) {
	client1, _ := NewClient("name", "email")
	account1 := NewAccount(client1)
	client1.AddAccount(account1)
	account1.Credit(100.0)

	client2, _ := NewClient("name", "email")
	account2 := NewAccount(client2)
	client2.AddAccount(account2)
	account2.Credit(100.0)

	transaction, err := NewTransaction(account1, account2, 10.0)

	assert.Nil(t, err)
	assert.NotNil(t, transaction)
	assert.NotEmpty(t, transaction.ID)
	assert.Equal(t, account1, transaction.AccountFrom)
	assert.Equal(t, account2, transaction.AccountTo)
	assert.Equal(t, 10.0, transaction.Amount)
	assert.Equal(t, 90.0, transaction.AccountFrom.Balance)
	assert.Equal(t, 90.0, account1.Balance)
	assert.Equal(t, 110.0, transaction.AccountTo.Balance)
	assert.Equal(t, 110.0, account2.Balance)
}


func TestCreateNewTransactionWithInsufficientBalance(t *testing.T) {
	client1, _ := NewClient("name", "email")
	account1 := NewAccount(client1)
	client1.AddAccount(account1)
	account1.Credit(100.0)

	client2, _ := NewClient("name", "email")
	account2 := NewAccount(client2)
	client2.AddAccount(account2)
	account2.Credit(100.0)

	transaction, err := NewTransaction(account1, account2, 200.0)

	assert.NotNil(t, err)
	assert.Error(t, err, "insufficient balance")
	assert.Nil(t, transaction)
	assert.Equal(t, 100.0, account1.Balance)
	assert.Equal(t, 100.0, account2.Balance)
}

func TestCreateNewTransactionWithSameAccounts(t *testing.T) {
	client1, _ := NewClient("name", "email")
	account1 := NewAccount(client1)
	client1.AddAccount(account1)
	account1.Credit(100.0)

	transaction, err := NewTransaction(account1, account1, 10.0)

	assert.NotNil(t, err)
	assert.Error(t, err, "accounts must be different")
	assert.Nil(t, transaction)
	assert.Equal(t, 100.0, account1.Balance)
}

func TestCreateNewTransactionWithInvalidAmount(t *testing.T) {
	client1, _ := NewClient("name", "email")
	account1 := NewAccount(client1)
	client1.AddAccount(account1)
	account1.Credit(100.0)
	
	client2, _ := NewClient("name", "email")
	account2 := NewAccount(client2)
	client2.AddAccount(account2)
	account2.Credit(100.0)

	transaction, err := NewTransaction(account1, account2, -10.0)

	assert.NotNil(t, err)
	assert.Error(t, err, "amount must be greater than zero")
	assert.Nil(t, transaction)
	assert.Equal(t, 100.0, account1.Balance)
	assert.Equal(t, 100.0, account2.Balance)
}
	