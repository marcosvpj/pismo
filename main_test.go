package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/marcosvpj/pismo/account"
	"github.com/stretchr/testify/assert"
)

func TestGetIndex(t *testing.T) {

	repository := account.NewMemoryRepository()
	service := account.NewService(repository)
	server := httptest.NewServer(NewAPIServer(service).Routes())
	defer server.Close()

	resp, err := http.Get(server.URL + "/")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
