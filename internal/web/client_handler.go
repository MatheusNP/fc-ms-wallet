package web

import (
	"encoding/json"
	"net/http"

	createclient "github.com/MatheusNP/fc-ms-wallet/internal/usecase/create_client"
)

type WebClientHandler struct {
	CreateClientUseCase createclient.CreateClientUseCase
}

func NewWebClientHandler(clientUseCase createclient.CreateClientUseCase) *WebClientHandler {
	return &WebClientHandler{
		CreateClientUseCase: clientUseCase,
	}
}

func (h *WebClientHandler) CreateClient(w http.ResponseWriter, r *http.Request) {
	var dto createclient.CreateClientInputDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte(err.Error()))
		return
	}

	output, err := h.CreateClientUseCase.Execute(dto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}
