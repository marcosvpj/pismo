package account

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAccountRepository struct {
	mock.Mock
}

func (m *MockAccountRepository) GetAccount(accountID int) (Account, error) {
	args := m.Called(accountID)
	return args.Get(0).(Account), args.Error(1)
}

func TestServiceGetAccountExpectAccount(t *testing.T) {
	mockRepository := new(MockAccountRepository)
	mockRepository.On("GetAccount", 1).Return(Account{AccountID: 1, DocumentNumber: "12345678900"}, nil)

	service := NewService(mockRepository)

	acc, err := service.GetAccount(1)

	expectedAcc := Account{
		AccountID:      1,
		DocumentNumber: "12345678900",
	}

	mockRepository.AssertExpectations(t)

	assert.Equal(t, expectedAcc, acc)
	assert.NoError(t, err)
}

func TestServiceGetAccountExpectNotFound(t *testing.T) {
	mockRepository := new(MockAccountRepository)
	mockRepository.On("GetAccount", 1).Return(Account{}, ErrNotFound)

	service := NewService(mockRepository)

	acc, err := service.GetAccount(1)

	expectedAcc := Account{}

	mockRepository.AssertExpectations(t)

	assert.Equal(t, expectedAcc, acc)
	assert.ErrorIs(t, err, ErrNotFound)
}
