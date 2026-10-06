package transaction

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) Save(ctx context.Context, t Transaction) (Transaction, error) {
	args := m.Called(t)
	return args.Get(0).(Transaction), args.Error(1)
}
