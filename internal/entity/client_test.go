package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateNewClient(t *testing.T) {
	client, err := NewClient("name", "email")

	assert.Nil(t, err)
	assert.NotNil(t, client)
	assert.NotEmpty(t, client.ID)
	assert.Equal(t, "name", client.Name)
	assert.Equal(t, "email", client.Email)
}

func TestCreateWhenArgumentsAreInvalid(t *testing.T) {
	_, err := NewClient("", "")

	assert.NotNil(t, err)
}

func TestUpdateClient(t *testing.T) {
	client, _ := NewClient("name", "email")
	err := client.Update("new name", "new email")

	assert.Nil(t, err)
	assert.Equal(t, "new name", client.Name)
	assert.Equal(t, "new email", client.Email)
}

func TestUpdateWhenArgumentsAreInvalid(t *testing.T) {
	client, _ := NewClient("name", "email")
	err := client.Update("", "")

	assert.Error(t, err, "name is required")
}
