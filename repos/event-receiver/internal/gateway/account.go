package gateway

import "github.com/MatheusNP/fc-ms-wallet/ms-account/internal/entity"

type AccountGateway interface {
	Save(account *entity.Account) error
	FindByID(id string) (*entity.Account, error)
	UpdateBalance(account *entity.Account) error
}
