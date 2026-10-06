package transaction

import (
	"testing"
	"time"

	"github.com/marcosvpj/pismo/account"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestServiceCreateTransactionExpectSuccess(t *testing.T) {
	now := func() time.Time { return time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC) }
	eventDate := now()
	// eventDate := func() time.Time { return time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC) }

	mockAccountRepository := new(account.MockAccountRepository)
	mockAccountRepository.On("FindByID", 1).Return(account.Account{AccountID: 1, DocumentNumber: "12345678900"}, nil)
	mockRepository := new(MockTransactionRepository)
	mockRepository.On("Save", Transaction{AccountID: 1, OperationTypeID: CreditVoucher, Amount: decimal.RequireFromString("123.45"), EventDate: eventDate}).Return(Transaction{TransactionID: 1, AccountID: 1, OperationTypeID: CreditVoucher, Amount: decimal.RequireFromString("123.45"), EventDate: eventDate}, nil)

	transaction, err := NewTransaction(1, 4, decimal.RequireFromString("123.45"))
	if err != nil {
		t.Fatal(err)
	}

	service := NewService(mockRepository, mockAccountRepository, now)
	tr, err := service.CreateTransaction(t.Context(), transaction)

	expectedTr := Transaction{
		TransactionID:   1,
		AccountID:       1,
		OperationTypeID: CreditVoucher,
		Amount:          decimal.RequireFromString("123.45"),
		EventDate:       eventDate,
	}

	mockRepository.AssertExpectations(t)

	assert.True(t, expectedTr.Amount.Equal(tr.Amount), "expected %s, got %s", expectedTr.Amount, tr.Amount)
	assert.Equal(t, 1, tr.AccountID)
	assert.Equal(t, OperationType(expectedTr.OperationTypeID), tr.OperationTypeID)
	assert.Equal(t, eventDate, tr.EventDate)
	assert.NoError(t, err)
}

// func TestServiceGetAccountExpectNotFound(t *testing.T) {
// 	mockRepository := new(MockAccountRepository)
// 	mockRepository.On("FindByID", 1).Return(Account{}, ErrNotFound)

// 	service := NewService(mockRepository)

// 	acc, err := service.GetAccount(t.Context(), 1)

// 	expectedAcc := Account{}

// 	mockRepository.AssertExpectations(t)

// 	assert.Equal(t, expectedAcc, acc)
// 	assert.ErrorIs(t, err, ErrNotFound)
// }

// func TestCreateAccountExpectOk(t *testing.T) {
// 	mockRepository := new(MockAccountRepository)
// 	mockRepository.On("Save", Account{DocumentNumber: "12345678900"}).Return(Account{AccountID: 1, DocumentNumber: "12345678900"}, nil)

// 	service := NewService(mockRepository)

// 	acc, err := service.CreateAccount(t.Context(), "12345678900")

// 	expectedAcc := Account{AccountID: 1, DocumentNumber: "12345678900"}

// 	mockRepository.AssertExpectations(t)

// 	assert.Equal(t, expectedAcc, acc)
// 	assert.NoError(t, err)
// }

// func TestCreateAccountExpectMissingField(t *testing.T) {
// 	mockRepository := new(MockAccountRepository)

// 	service := NewService(mockRepository)

// 	acc, err := service.CreateAccount(t.Context(), "")

// 	expectedAcc := Account{}

// 	mockRepository.AssertNotCalled(t, "Save", mock.Anything)

// 	assert.Equal(t, expectedAcc, acc)
// 	assert.ErrorIs(t, err, ErrFieldDocumentNumberMissing)
// }
