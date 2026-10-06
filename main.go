package main

import (
	"os"
	"time"

	"github.com/marcosvpj/pismo/account"
	"github.com/marcosvpj/pismo/transaction"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	newApp(func() time.Time { return time.Now().UTC() }).Serve(":" + port)
}

func newApp(now func() time.Time) *APIServer {
	accountRepository := account.NewMemoryRepository()
	transactionRepository := transaction.NewMemoryRepository()

	return NewAPIServer(
		account.NewService(accountRepository),
		transaction.NewService(transactionRepository, accountRepository, now),
	)
}
