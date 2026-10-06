package transaction

import (
	"errors"
	"time"

	"github.com/shopspring/decimal"
)

type Transaction struct {
	TransactionID   int             `json:"transaction_id"`
	AccountID       int             `json:"account_id"`
	OperationTypeID OperationType   `json:"operation_type_id"`
	Amount          decimal.Decimal `json:"amount"`
	EventDate       time.Time       `json:"event_date"`
}

var ErrInvalidOperationType = errors.New("invalid operation type")
var ErrInvalidAmount = errors.New("invalid amount")
var ErrInvalidDecimalSize = errors.New("invalid decimal size")

type OperationType int

const (
	OperationInvalid OperationType = iota
	NormalPurchase
	PurchaseWithInstallments
	Withdrawal
	CreditVoucher
	operationCount
)

var operationNames = map[OperationType]string{
	OperationInvalid:         "Operation invalid",
	NormalPurchase:           "Normal purchase",
	PurchaseWithInstallments: "Purchase with installments",
	Withdrawal:               "Withdrawal",
	CreditVoucher:            "Credit voucher",
}

func (o OperationType) String() string {
	if name, ok := operationNames[o]; ok {
		return name
	}
	return OperationInvalid.String()
}

func (o OperationType) IsValid() bool {
	return o > OperationInvalid && o < operationCount
}

func (o OperationType) IsDebit() bool {
	switch o {
	case NormalPurchase, PurchaseWithInstallments, Withdrawal:
		return true
	case CreditVoucher:
		return false
	default:
		return false
	}
}

func NewTransaction(accountID int, opType OperationType, amount decimal.Decimal) (Transaction, error) {
	if !amount.Equal(amount.Round(2)) {
		return Transaction{}, ErrInvalidDecimalSize
	}

	if !opType.IsValid() {
		return Transaction{}, ErrInvalidOperationType
	}

	if amount.IsZero() {
		return Transaction{}, ErrInvalidAmount
	}

	amount = amount.Abs()
	if opType.IsDebit() {
		amount = amount.Neg()
	}

	return Transaction{
		AccountID:       accountID,
		OperationTypeID: opType,
		Amount:          amount,
	}, nil
}
