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
	server := httptest.NewServer(newApp(fixedNow).Routes())
	defer server.Close()

	// create account
	expected := account.Account{
		AccountID:      1,
		DocumentNumber: "12345678900",
	}

	payload := bytes.NewBufferString(`{"document_number":"12345678900"}`)

	resp, err := http.Post(server.URL+"/accounts", "application/json", payload)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	response, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	expectedBody := `{"account_id":1, "document_number":"12345678900"}`
	assert.JSONEq(t, expectedBody, string(response))

	// get account
	resp, err = http.Get(server.URL + "/accounts/1")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	var acc account.Account
	err = json.NewDecoder(resp.Body).Decode(&acc)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, expected, acc)

	// create negative transaction
	payloadTransaction := bytes.NewBufferString(`{"account_id": 1,"operation_type_id": 3,"amount": -123.45}`)

	resp, err = http.Post(server.URL+"/transactions", "application/json", payloadTransaction)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	response, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	expectedTransactionBody := `{"transaction_id":1,"account_id": 1,"operation_type_id": 3,"amount": -123.45, "event_date":"2009-11-17T20:34:58.651387237Z"}`
	assert.JSONEq(t, expectedTransactionBody, string(response))

	// create positive transaction
	payloadPositiveTransaction := bytes.NewBufferString(`{"account_id": 1,"operation_type_id": 4,"amount": 123.45}`)

	resp, err = http.Post(server.URL+"/transactions", "application/json", payloadPositiveTransaction)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	response, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	expectedPositiveTransactionBody := `{"transaction_id":2,"account_id": 1,"operation_type_id": 4,"amount": 123.45, "event_date":"2009-11-17T20:34:58.651387237Z"}`
	assert.JSONEq(t, expectedPositiveTransactionBody, string(response))
}
