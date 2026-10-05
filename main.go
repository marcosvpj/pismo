package main

import (
	"github.com/marcosvpj/pismo/account"
)

func main() {
	repository := account.NewMemoryRepository()
	service := account.NewService(repository)
	NewAPIServer(service).Serve()
}
