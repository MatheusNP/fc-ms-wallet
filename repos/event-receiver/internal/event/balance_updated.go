package event

import "time"

type BalanceUpdated struct {
	Name    string      `json:"name"`
	Payload interface{} `json:"payload"`
}

func NewBalanceUpdated() *BalanceUpdated {
	return &BalanceUpdated{
		Name: "BalanceUpdated",
	}
}

func (b *BalanceUpdated) GetName() string {
	return b.Name
}

func (b *BalanceUpdated) GetPayload() interface{} {
	return b.Payload
}

func (b *BalanceUpdated) SetPayload(payload interface{}) {
	b.Payload = payload
}

func (b *BalanceUpdated) GetDateTime() time.Time {
	return time.Now()
}
