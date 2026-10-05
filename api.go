package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/marcosvpj/pismo/account"
)

type APIServer struct {
	accountService *account.Service
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

func NewAPIServer(accountService *account.Service) *APIServer {
	return &APIServer{
		accountService: accountService,
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
		writeError(w, http.StatusBadRequest, "invalid account information")
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

func (a *APIServer) getHealthHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (a *APIServer) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /accounts/{account_id}", a.getAccountHandler)
	mux.HandleFunc("POST /accounts", a.postAccountHandler)
	mux.HandleFunc("GET /health", a.getHealthHandler)
	return mux
}

func (a *APIServer) Serve() {
	server := &http.Server{
		Addr:         ":8080",
		Handler:      a.Routes(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}
