package transaction

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidOperationTypeIDIsValid(t *testing.T) {
	tests := []struct {
		name            string
		operationTypeId int
		expected        bool
		expectedName    string
	}{
		{
			name:            "Normal purchase",
			expectedName:    "Normal purchase",
			operationTypeId: 1,
			expected:        true,
		}, {
			name:            "Purchase with installments",
			expectedName:    "Purchase with installments",
			operationTypeId: 2,
			expected:        true,
		}, {
			name:            "Withdrawal",
			expectedName:    "Withdrawal",
			operationTypeId: 3,
			expected:        true,
		}, {
			name:            "Credit voucher",
			expectedName:    "Credit voucher",
			operationTypeId: 4,
			expected:        true,
		}, {
			name:            "Invalid operation ID",
			expectedName:    "Operation invalid",
			operationTypeId: 5,
			expected:        false,
		}, {
			name:            "Invalid operation ID",
			expectedName:    "Operation invalid",
			operationTypeId: 0,
			expected:        false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, OperationType(test.operationTypeId).IsValid())
			assert.Equal(t, test.expectedName, OperationType(test.operationTypeId).String())
		})
	}
}

func TestNewTransaction(t *testing.T) {
	tests := []struct {
		name            string
		operationTypeID int
		amount          string
		expectedAmount  string
		expectedError   error
	}{
		{name: "Valid normal purchase", operationTypeID: 1, amount: "-100", expectedAmount: "-100"},
		{name: "Positive value normal purchase", operationTypeID: 1, amount: "100", expectedAmount: "-100"},
		{name: "Not valid transaction", operationTypeID: 10, amount: "-100", expectedError: ErrInvalidOperationType},
		{name: "Normal purchase is negative", operationTypeID: 1, amount: "50.00", expectedAmount: "-50"},
		{name: "Purchase with installments is negative", operationTypeID: 2, amount: "23.5", expectedAmount: "-23.5"},
		{name: "Withdrawal is negative", operationTypeID: 3, amount: "18.7", expectedAmount: "-18.7"},
		{name: "Credit voucher is positive", operationTypeID: 4, amount: "60", expectedAmount: "60"},
		{name: "Zero amount", operationTypeID: 1, amount: "0", expectedError: ErrInvalidAmount},
		{name: "Purchase is negative", operationTypeID: 1, amount: "50", expectedAmount: "-50"},
		{name: "Negative credit voucher is positive", operationTypeID: 4, amount: "-60", expectedAmount: "60"},
		{name: "Unknown operation type", operationTypeID: 5, amount: "50", expectedError: ErrInvalidOperationType},
		{name: "Zero operation type", operationTypeID: 0, amount: "50", expectedError: ErrInvalidOperationType},
		{name: "Too many decimals", operationTypeID: 1, amount: "50.123", expectedError: ErrInvalidDecimalSize},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tr, err := NewTransaction(1, OperationType(test.operationTypeID), decimal.RequireFromString(test.amount))
			require.ErrorIs(t, err, test.expectedError)
			if err != nil {
				return
			}

			expectedAmount := decimal.RequireFromString(test.expectedAmount)
			assert.True(t, expectedAmount.Equal(tr.Amount), "expected %s, got %s", expectedAmount, tr.Amount)
			assert.Equal(t, 1, tr.AccountID)
			assert.Equal(t, OperationType(test.operationTypeID), tr.OperationTypeID)
		})
	}
}
