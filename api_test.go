package main

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/marcosvpj/pismo/account"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetHealth(t *testing.T) {
	repository := account.NewMemoryRepository()
	service := account.NewService(repository)
	server := httptest.NewServer(NewAPIServer(service).Routes())
	defer server.Close()

	resp, err := http.Get(server.URL + "/health")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, resp.StatusCode)
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
			mockRepository := new(account.MockAccountRepository)
			accountService := account.NewService(mockRepository)
			if test.setup != nil {
				test.setup(mockRepository)
			}
			api := NewAPIServer(accountService)

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
			mockRepository.AssertExpectations(t)
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
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockRepository := new(account.MockAccountRepository)
			accountService := account.NewService(mockRepository)
			if test.setup != nil {
				test.setup(mockRepository)
			}
			api := NewAPIServer(accountService)

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
			mockRepository.AssertExpectations(t)
		})
	}
}
