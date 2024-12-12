package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	getaccountbalance "github.com/MatheusNP/fc-ms-wallet/ms-account/internal/usecase/get_account_balance"
)

type WebAccountHandler struct {
	GetAccountBalanceUseCase getaccountbalance.GetAccountBalanceUseCase
}

func NewWebAccountHandler(getAccountBalanceUseCase getaccountbalance.GetAccountBalanceUseCase) *WebAccountHandler {
	return &WebAccountHandler{
		GetAccountBalanceUseCase: getAccountBalanceUseCase,
	}
}

func (h *WebAccountHandler) GetAccountBalance(w http.ResponseWriter, r *http.Request) {
	dto := getaccountbalance.GetAccountBalanceInputDTO{
		AccountID: r.URL.Query().Get("account_id"),
	}

	w.Header().Set("Content-Type", "application/json")

	output, err := h.GetAccountBalanceUseCase.Execute(dto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err.Error())
		return
	}

	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
}
