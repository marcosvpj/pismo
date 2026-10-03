package account

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestServiceGetAccountExpectAccount(t *testing.T) {
	mockRepository := new(MockAccountRepository)
	mockRepository.On("FindByID", 1).Return(Account{AccountID: 1, DocumentNumber: "12345678900"}, nil)

	service := NewService(mockRepository)

	acc, err := service.GetAccount(t.Context(), 1)

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
	mockRepository.On("FindByID", 1).Return(Account{}, ErrNotFound)

	service := NewService(mockRepository)

	acc, err := service.GetAccount(t.Context(), 1)

	expectedAcc := Account{}

	mockRepository.AssertExpectations(t)

	assert.Equal(t, expectedAcc, acc)
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestCreateAccountExpectOk(t *testing.T) {
	mockRepository := new(MockAccountRepository)
	mockRepository.On("SaveAccount", Account{DocumentNumber: "12345678900"}).Return(Account{AccountID: 1, DocumentNumber: "12345678900"}, nil)

	service := NewService(mockRepository)

	acc, err := service.CreateAccount(t.Context(), Account{DocumentNumber: "12345678900"})

	expectedAcc := Account{AccountID: 1, DocumentNumber: "12345678900"}

	mockRepository.AssertExpectations(t)

	assert.Equal(t, expectedAcc, acc)
	assert.NoError(t, err)
}

func TestCreateAccountExpectMissingField(t *testing.T) {
	mockRepository := new(MockAccountRepository)

	service := NewService(mockRepository)

	acc, err := service.CreateAccount(t.Context(), Account{DocumentNumber: ""})

	expectedAcc := Account{}

	mockRepository.AssertNotCalled(t, "SaveAccount", mock.Anything)

	assert.Equal(t, expectedAcc, acc)
	assert.ErrorIs(t, err, ErrFieldDocumentNumberMissing)
}
