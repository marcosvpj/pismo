package main

import (
	"time"

	"github.com/marcosvpj/pismo/account"
	"github.com/marcosvpj/pismo/transaction"
)

func main() {
	newApp(func() time.Time { return time.Now().UTC() }).Serve()
}

func newApp(now func() time.Time) *APIServer {
	accountRepository := account.NewMemoryRepository()
	transactionRepository := transaction.NewMemoryRepository()

	return NewAPIServer(
		account.NewService(accountRepository),
		transaction.NewService(transactionRepository, accountRepository, now),
	)
}
