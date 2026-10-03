package account

import (
	"github.com/stretchr/testify/mock"
)

type MockAccountRepository struct {
	mock.Mock
}

func (m *MockAccountRepository) FindByID(accountID int) (Account, error) {
	args := m.Called(accountID)
	return args.Get(0).(Account), args.Error(1)
}

func (m *MockAccountRepository) SaveAccount(acc Account) (Account, error) {
	args := m.Called(acc)
	return args.Get(0).(Account), args.Error(1)
}
