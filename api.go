package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/marcosvpj/pismo/account"
	"github.com/marcosvpj/pismo/transaction"
	"github.com/shopspring/decimal"
)

type APIServer struct {
	accountService     *account.Service
	transactionService *transaction.Service
}

type errorResponse struct {
	Error string `json:"error"`
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, errorResponse{Error: msg})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func NewAPIServer(accountService *account.Service, transactionService *transaction.Service) *APIServer {
	return &APIServer{
		accountService:     accountService,
		transactionService: transactionService,
	}
}

func (a *APIServer) getAccountHandler(w http.ResponseWriter, r *http.Request) {
	accountID, err := strconv.Atoi(r.PathValue("account_id"))
	if err != nil || accountID <= 0 {
		writeError(w, http.StatusBadRequest, "invalid account id")
		return
	}

	acc, err := a.accountService.GetAccount(r.Context(), accountID)
	if errors.Is(err, account.ErrNotFound) {
		writeError(w, http.StatusNotFound, err.Error())
		return
	} else if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJSON(w, http.StatusOK, acc)
}

type createAccountRequest struct {
	DocumentNumber string `json:"document_number"`
}

func (a *APIServer) postAccountHandler(w http.ResponseWriter, r *http.Request) {
	var acc createAccountRequest
	err := json.NewDecoder(r.Body).Decode(&acc)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	createdAccount, err := a.accountService.CreateAccount(r.Context(), acc.DocumentNumber)
	if errors.Is(err, account.ErrFieldDocumentNumberMissing) {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	} else if err != nil {
		writeError(w, http.StatusInternalServerError, "error creating account")
		return
	}

	writeJSON(w, http.StatusCreated, createdAccount)
}

type createTransactionRequest struct {
	AccountID       int             `json:"account_id"`
	OperationTypeID int             `json:"operation_type_id"`
	Amount          decimal.Decimal `json:"amount"`
}

func (a *APIServer) postTransactionHandler(w http.ResponseWriter, r *http.Request) {
	var tInput createTransactionRequest
	err := json.NewDecoder(r.Body).Decode(&tInput)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	if tInput.AccountID == 0 {
		writeError(w, http.StatusBadRequest, "account_id is required")
		return
	}
	if tInput.OperationTypeID == 0 {
		writeError(w, http.StatusBadRequest, "operation_type_id is required")
		return
	}
	if tInput.Amount.IsZero() {
		writeError(w, http.StatusBadRequest, "amount is required")
		return
	}

	newTransaction, err := transaction.NewTransaction(tInput.AccountID, transaction.OperationType(tInput.OperationTypeID), tInput.Amount)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid transaction data")
		return
	}

	createdTransaction, err := a.transactionService.CreateTransaction(r.Context(), newTransaction)
	if errors.Is(err, transaction.ErrInvalidAccount) {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	} else if err != nil {
		writeError(w, http.StatusInternalServerError, "error creating transaction")
		return
	}

	writeJSON(w, http.StatusCreated, createdTransaction)

}

func (a *APIServer) getHealthHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (a *APIServer) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /accounts/{account_id}", a.getAccountHandler)
	mux.HandleFunc("POST /accounts", a.postAccountHandler)
	mux.HandleFunc("POST /transactions", a.postTransactionHandler)
	mux.HandleFunc("GET /health", a.getHealthHandler)
	return mux
}

func (a *APIServer) Serve() {
	log.Println("Starting server on port :8080")

	server := &http.Server{
		Addr:         ":8080",
		Handler:      a.Routes(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}

func init() {
	decimal.MarshalJSONWithoutQuotes = true
}
