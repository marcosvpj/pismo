package account

import "errors"

type Repository interface {
	GetAccount(accountID int) (Account, error)
}

var ErrNotFound = errors.New("account not found")
