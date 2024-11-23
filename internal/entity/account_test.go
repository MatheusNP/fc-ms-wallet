package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateNewAccount(t *testing.T) {
	client, _ := NewClient("name", "email")
	account := NewAccount(client)

	assert.NotNil(t, account)
	assert.NotEmpty(t, account.ID)
	assert.Equal(t, client, account.Client)
	assert.Equal(t, 0.0, account.Balance)
}

func TestCreateAccountWithNilClient(t *testing.T){
	account := NewAccount(nil)

	assert.Nil(t, account)
}

func TestCredit(t *testing.T) {
	client, _ := NewClient("name", "email")
	account := NewAccount(client)

	account.Credit(10.0)
	assert.Equal(t, 10.0, account.Balance)

	account.Credit(10.0)	
	assert.Equal(t, 20.0, account.Balance)
}

func TestDebit(t *testing.T) {
	client, _ := NewClient("name", "email")
	account := NewAccount(client)

	account.Credit(20.0)
	account.Debit(5.0)
	assert.Equal(t, 15.0, account.Balance)
}