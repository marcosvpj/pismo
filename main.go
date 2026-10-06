package main

import (
	"time"

	"github.com/marcosvpj/pismo/account"
	"github.com/marcosvpj/pismo/transaction"
)

func main() {
	now := func() time.Time { return time.Now().UTC() }

	accountRepository := account.NewMemoryRepository()
	accountService := account.NewService(accountRepository)

	transactionRepository := transaction.NewMemoryRepository()
	transactionService := transaction.NewService(transactionRepository, accountRepository, now)

	NewAPIServer(accountService, transactionService).Serve()
}
