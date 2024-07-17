package web

import (
	"encoding/json"
	"net/http"

	createtransaction "github.com.br/derivedpuma7/wallet-core/internal/usecase/create_transaction"
)

type WebTransactionHandler struct {
  CreateTransacitonUseCase createtransaction.CreateTransactionUseCase
}

func NewWebTransactionHandler(createTransactionUseCase createtransaction.CreateTransactionUseCase) *WebTransactionHandler {
  return &WebTransactionHandler{
    CreateTransacitonUseCase: createTransactionUseCase,
  }
}

func (h *WebTransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
  var dto createtransaction.CreateTransactionInputDto
  err := json.NewDecoder(r.Body).Decode(&dto)
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    return
  }
  
  output, err := h.CreateTransacitonUseCase.Execute(r.Context(), dto)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    w.Write([]byte(err.Error()))
    return
  }

  w.Header().Set("Content-Type", "application/json")
  err = json.NewEncoder(w).Encode(output)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  w.WriteHeader(http.StatusCreated)
}
