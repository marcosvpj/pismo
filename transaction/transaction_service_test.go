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

	mockAccountRepository := new(account.MockAccountRepository)
	mockAccountRepository.On("FindByID", 1).Return(account.Account{AccountID: 1, DocumentNumber: "12345678900"}, nil)
	mockRepository := new(MockTransactionRepository)
	mockRepository.On("Save", Transaction{AccountID: 1, OperationTypeID: CreditVoucher, Amount: decimal.RequireFromString("123.45"), EventDate: now()}).Return(Transaction{TransactionID: 1, AccountID: 1, OperationTypeID: CreditVoucher, Amount: decimal.RequireFromString("123.45"), EventDate: now()}, nil)

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
		EventDate:       now(),
	}

	mockRepository.AssertExpectations(t)

	assert.True(t, expectedTr.Amount.Equal(tr.Amount), "expected %s, got %s", expectedTr.Amount, tr.Amount)
	assert.Equal(t, 1, tr.AccountID)
	assert.Equal(t, OperationType(expectedTr.OperationTypeID), tr.OperationTypeID)
	assert.Equal(t, now(), tr.EventDate)
	assert.NoError(t, err)
}
