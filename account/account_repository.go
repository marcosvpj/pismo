package account

import "errors"

type Repository interface {
	FindByID(accountID int) (Account, error)
	SaveAccount(account Account) (Account, error)
}

var ErrNotFound = errors.New("account not found")

type DB struct {
}

func (db *DB) FindByID(accountID int) (Account, error) {

	return Account{}, nil
}
func (db *DB) SaveAccount(account Account) (Account, error) {

	return Account{}, nil
}
