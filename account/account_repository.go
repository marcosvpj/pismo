package account

import (
	"context"
	"errors"
)

type Repository interface {
	FindByID(ctx context.Context, accountID int) (Account, error)
	Save(ctx context.Context, account Account) (Account, error)
}

var ErrNotFound = errors.New("account not found")
