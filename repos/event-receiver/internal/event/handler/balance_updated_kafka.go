package handler

import (
	"encoding/json"
	"fmt"
	"sync"

	updateaccountbalance "github.com/MatheusNP/fc-ms-wallet/ms-account/internal/usecase/update_account_balance"
	"github.com/MatheusNP/fc-ms-wallet/ms-account/pkg/events"
)

type UpdateBalanceKafkaHandler struct {
	UpdateBalance *updateaccountbalance.UpdateAccountBalanceUseCase
}

type BalanceUpdateKafkaInputDTO struct {
	AccountFromID      string  `json:"account_from_id"`
	AccountToID        string  `json:"account_to_id"`
	BalanceAccountFrom float64 `json:"balance_account_from"`
	BalanceAccountTo   float64 `json:"balance_account_to"`
}

func NewUpdateBalanceKafkaHandler(updateBalance *updateaccountbalance.UpdateAccountBalanceUseCase) *UpdateBalanceKafkaHandler {
	return &UpdateBalanceKafkaHandler{UpdateBalance: updateBalance}
}

func (h *UpdateBalanceKafkaHandler) Handle(msg events.EventInterface, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("UpdateBalanceKafkaHandler: ", msg.GetPayload())

	payload, err := json.Marshal(msg.GetPayload())
	if err != nil {
		panic(fmt.Sprintf("Erro ao converter o payload para string: %v", msg.GetPayload()))
	}

	var input BalanceUpdateKafkaInputDTO

	err = json.Unmarshal([]byte(payload), &input)
	if err != nil {
		panic(fmt.Sprintf("Erro ao decodificar JSON: %v", err))
	}

	err = h.UpdateBalance.Execute(updateaccountbalance.UpdateAccountBalanceInputDTO{
		AccountID: input.AccountFromID,
		Balance:   input.BalanceAccountFrom,
	})
	if err != nil {
		panic(fmt.Sprintf("Erro ao atualizar saldo: %v", err))
	}

	err = h.UpdateBalance.Execute(updateaccountbalance.UpdateAccountBalanceInputDTO{
		AccountID: input.AccountToID,
		Balance:   input.BalanceAccountTo,
	})
	if err != nil {
		panic(fmt.Sprintf("Erro ao atualizar saldo: %v", err))
	}

}
