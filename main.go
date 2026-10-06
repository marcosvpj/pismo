package main

import (
	"github.com/marcosvpj/pismo/account"
	"github.com/marcosvpj/pismo/transaction"
)

func main() {
	accountRepository := account.NewMemoryRepository()
	accountService := account.NewService(accountRepository)

	transactionRepository := transaction.NewMemoryRepository()
	transactionService := transaction.NewService(transactionRepository, accountRepository)

	NewAPIServer(accountService, transactionService).Serve()
}
