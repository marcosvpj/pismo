package transaction

import (
	"context"
	"errors"
)

type Repository interface {
	Save(ctx context.Context, transaction Transaction) (Transaction, error)
}

var ErrNotFound = errors.New("transaction not found")
