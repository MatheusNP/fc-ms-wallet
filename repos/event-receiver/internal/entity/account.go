package entity

import (
	"time"
)

type Account struct {
	ID        string
	Balance   float64
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewAccount(id string, balance float64) *Account {
	return &Account{
		ID:        id,
		Balance:   balance,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
