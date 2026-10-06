package main

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/marcosvpj/pismo/account"
	"github.com/marcosvpj/pismo/transaction"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetHealth(t *testing.T) {
	now := func() time.Time { return time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC) }
	mockAccountRepository := new(account.MockAccountRepository)
	mockTransactionRepository := new(transaction.MockTransactionRepository)
	accountService := account.NewService(mockAccountRepository)
	transactionService := transaction.NewService(mockTransactionRepository, mockAccountRepository, now)
	server := httptest.NewServer(NewAPIServer(accountService, transactionService).Routes())
	defer server.Close()

	res, err := http.Get(server.URL + "/health")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.JSONEq(t, `{"status":"ok"}`, string(body))
}

func TestGetAccountEndpoint(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		setup          func(m *account.MockAccountRepository)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Found",
			path: "/accounts/1",
			setup: func(m *account.MockAccountRepository) {
				m.On("FindByID", 1).Return(account.Account{AccountID: 1, DocumentNumber: "12345678900"}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"account_id": 1, "document_number": "12345678900"}`,
		},
		{
			name: "Not found",
			path: "/accounts/1",
			setup: func(m *account.MockAccountRepository) {
				m.On("FindByID", 1).Return(account.Account{}, account.ErrNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":"account not found"}`,
		},
		{
			name: "Internal error",
			path: "/accounts/1",
			setup: func(m *account.MockAccountRepository) {
				m.On("FindByID", 1).Return(account.Account{}, errors.New("db down"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"internal error"}`,
		},
		{
			name:           "Invalid account id",
			path:           "/accounts/xxx",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"invalid account id"}`,
		},
		{
			name:           "Zero account id",
			path:           "/accounts/0",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"invalid account id"}`,
		},
		{
			name:           "Negative account id",
			path:           "/accounts/-1",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"invalid account id"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			now := func() time.Time { return time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC) }
			mockAccountRepository := new(account.MockAccountRepository)
			mockTransactionRepository := new(transaction.MockTransactionRepository)
			accountService := account.NewService(mockAccountRepository)
			transactionService := transaction.NewService(mockTransactionRepository, mockAccountRepository, now)
			if test.setup != nil {
				test.setup(mockAccountRepository)
			}
			api := NewAPIServer(accountService, transactionService)

			req := httptest.NewRequest(http.MethodGet, test.path, nil)
			w := httptest.NewRecorder()

			api.Routes().ServeHTTP(w, req)

			res := w.Result()
			defer res.Body.Close()
			body, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			assert.Equal(t, test.expectedStatus, res.StatusCode)
			assert.Equal(t, "application/json", res.Header.Get("Content-Type"))
			assert.JSONEq(t, test.expectedBody, string(body))
			mockAccountRepository.AssertExpectations(t)
		})
	}
}

func TestCreateAccountEndpoint(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		setup          func(m *account.MockAccountRepository)
		payload        string
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Create account",
			path: "/accounts",
			setup: func(m *account.MockAccountRepository) {
				m.On("Save", account.Account{DocumentNumber: "12345678900"}).Return(account.Account{AccountID: 1, DocumentNumber: "12345678900"}, nil)
			},
			payload:        `{"document_number":"12345678900"}`,
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"account_id": 1, "document_number": "12345678900"}`,
		},
		{
			name:           "Error missing field",
			path:           "/accounts",
			payload:        `{"document_number":""}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error": "field document_number is required"}`,
		},
		{
			name: "Error 500",
			path: "/accounts",
			setup: func(m *account.MockAccountRepository) {
				m.On("Save", account.Account{DocumentNumber: "12345678900"}).Return(account.Account{}, errors.New("db down"))
			},
			payload:        `{"document_number":"12345678900"}`,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error": "error creating account"}`,
		},
		{
			name:           "Error invalid JSON",
			path:           "/accounts",
			payload:        `{invalid`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error": "invalid json"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			now := func() time.Time { return time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC) }
			mockAccountRepository := new(account.MockAccountRepository)
			mockTransactionRepository := new(transaction.MockTransactionRepository)
			accountService := account.NewService(mockAccountRepository)
			transactionService := transaction.NewService(mockTransactionRepository, mockAccountRepository, now)
			if test.setup != nil {
				test.setup(mockAccountRepository)
			}
			api := NewAPIServer(accountService, transactionService)

			req := httptest.NewRequest(http.MethodPost, test.path, bytes.NewBufferString(test.payload))
			w := httptest.NewRecorder()

			api.Routes().ServeHTTP(w, req)

			res := w.Result()
			defer res.Body.Close()
			body, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			assert.Equal(t, test.expectedStatus, res.StatusCode)
			assert.Equal(t, "application/json", res.Header.Get("Content-Type"))
			assert.JSONEq(t, test.expectedBody, string(body))
			mockAccountRepository.AssertExpectations(t)
		})
	}
}

func TestCreateTransactionEndpoint(t *testing.T) {
	now := func() time.Time { return time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC) }

	tests := []struct {
		name           string
		path           string
		setup          func(ma *account.MockAccountRepository, mt *transaction.MockTransactionRepository)
		payload        string
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Create transaction with positive amount",
			path: "/transactions",
			setup: func(ma *account.MockAccountRepository, mt *transaction.MockTransactionRepository) {
				ma.On("FindByID", 1).Return(account.Account{AccountID: 1, DocumentNumber: "12345678900"}, nil)
				mt.On("Save", transaction.Transaction{AccountID: 1, OperationTypeID: transaction.CreditVoucher, Amount: decimal.RequireFromString("123.45"), EventDate: now()}).Return(transaction.Transaction{TransactionID: 1, AccountID: 1, OperationTypeID: transaction.CreditVoucher, Amount: decimal.RequireFromString("123.45"), EventDate: now()}, nil)
			},
			payload:        `{"account_id": 1,"operation_type_id": 4,"amount": 123.45}`,
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"transaction_id":1,"account_id": 1,"operation_type_id": 4,"amount": 123.45, "event_date":"2009-11-17T20:34:58.651387237Z"}`,
		},
		{
			name: "Create transaction with negative amount",
			path: "/transactions",
			setup: func(ma *account.MockAccountRepository, mt *transaction.MockTransactionRepository) {
				ma.On("FindByID", 1).Return(account.Account{AccountID: 1, DocumentNumber: "12345678900"}, nil)
				mt.On("Save", transaction.Transaction{AccountID: 1, OperationTypeID: transaction.Withdrawal, Amount: decimal.RequireFromString("-123.45"), EventDate: now()}).Return(transaction.Transaction{TransactionID: 1, AccountID: 1, OperationTypeID: transaction.Withdrawal, Amount: decimal.RequireFromString("-123.45"), EventDate: now()}, nil)
			},
			payload:        `{"account_id": 1,"operation_type_id": 3,"amount": -123.45}`,
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"transaction_id":1,"account_id": 1,"operation_type_id": 3,"amount": -123.45, "event_date":"2009-11-17T20:34:58.651387237Z"}`,
		},
		{
			name:           "Create transaction with missing account_id",
			path:           "/transactions",
			payload:        `{"operation_type_id": 3,"amount": -123.45}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error": "account_id is required"}`,
		},
		{
			name:           "Create transaction with missing operation_type_id",
			path:           "/transactions",
			payload:        `{"account_id": 1,"amount": -123.45}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error": "operation_type_id is required"}`,
		},
		{
			name:           "Create transaction with missing amount",
			path:           "/transactions",
			payload:        `{"account_id": 1,"operation_type_id": 3}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error": "amount is required"}`,
		},
		{
			name:           "Create transaction with invalid JSON",
			path:           "/transactions",
			payload:        `{"invalid`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error": "invalid json"}`,
		},
		{
			name: "Create transaction with error 500",
			path: "/transactions",
			setup: func(ma *account.MockAccountRepository, mt *transaction.MockTransactionRepository) {
				ma.On("FindByID", 1).Return(account.Account{AccountID: 1, DocumentNumber: "12345678900"}, nil)
				mt.On("Save", transaction.Transaction{AccountID: 1, OperationTypeID: transaction.Withdrawal, Amount: decimal.RequireFromString("-123.45"), EventDate: now()}).Return(transaction.Transaction{}, errors.New("db down"))
			},
			payload:        `{"account_id": 1,"operation_type_id": 3,"amount": -123.45}`,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error": "error creating transaction"}`,
		},
		{
			name: "Create transaction with invalid account",
			path: "/transactions",
			setup: func(ma *account.MockAccountRepository, mt *transaction.MockTransactionRepository) {
				ma.On("FindByID", 1).Return(account.Account{}, account.ErrNotFound)
				mt.On("Save", transaction.Transaction{AccountID: 1, OperationTypeID: transaction.Withdrawal, Amount: decimal.RequireFromString("-123.45"), EventDate: now()}).Return(transaction.Transaction{}, errors.New("db down"))
			},
			payload:        `{"account_id": 1,"operation_type_id": 3,"amount": -123.45}`,
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   `{"error": "invalid account"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockAccountRepository := new(account.MockAccountRepository)
			mockTransactionRepository := new(transaction.MockTransactionRepository)
			accountService := account.NewService(mockAccountRepository)
			transactionService := transaction.NewService(mockTransactionRepository, mockAccountRepository, now)
			if test.setup != nil {
				test.setup(mockAccountRepository, mockTransactionRepository)
			}
			api := NewAPIServer(accountService, transactionService)

			req := httptest.NewRequest(http.MethodPost, test.path, bytes.NewBufferString(test.payload))
			w := httptest.NewRecorder()

			api.Routes().ServeHTTP(w, req)

			res := w.Result()
			defer res.Body.Close()
			body, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			assert.Equal(t, test.expectedStatus, res.StatusCode)
			assert.Equal(t, "application/json", res.Header.Get("Content-Type"))
			assert.JSONEq(t, test.expectedBody, string(body))
			mockAccountRepository.AssertExpectations(t)
		})
	}
}
