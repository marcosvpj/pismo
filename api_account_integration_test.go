package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/marcosvpj/pismo/account"
	"github.com/stretchr/testify/assert"
)

func TestCreateAndRetrieveAccount(t *testing.T) {
	repository := account.NewMemoryRepository()
	service := account.NewService(repository)
	server := httptest.NewServer(NewAPIServer(service).Routes())
	defer server.Close()

	expected := account.Account{
		AccountID:      1,
		DocumentNumber: "12345678900",
	}

	payload := bytes.NewBufferString(`{"document_number":"12345678900"}`)

	resp, err := http.Post(server.URL+"/accounts", "application/json", payload)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	response, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		t.Fatal(err)
	}

	expectedBody := `{"account_id":1, "document_number":"12345678900"}`
	assert.JSONEq(t, expectedBody, string(response))

	resp, err = http.Get(server.URL + "/accounts/1")
	if err != nil {
		t.Fatal(err)
	}

	var acc account.Account
	err = json.NewDecoder(resp.Body).Decode(&acc)
	defer resp.Body.Close()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, expected, acc)
}
