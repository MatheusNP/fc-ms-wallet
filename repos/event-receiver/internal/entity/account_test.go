package entity

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateNewAccount(t *testing.T) {
	account := NewAccount(uuid.NewString(), 0.0)

	assert.NotNil(t, account)
	assert.NotEmpty(t, account.ID)
	assert.Equal(t, 0.0, account.Balance)
}
