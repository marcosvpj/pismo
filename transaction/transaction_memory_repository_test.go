package transaction

import (
	"sync"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestTransaction() Transaction {
	return Transaction{
		AccountID:       1,
		OperationTypeID: Withdrawal,
		Amount:          decimal.RequireFromString("-50.25"),
		EventDate:       time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC),
	}
}

func TestMemoryRepositorySaveAssignsSequentialIDs(t *testing.T) {
	repository := NewMemoryRepository()

	first, err := repository.Save(t.Context(), newTestTransaction())
	require.NoError(t, err)
	second, err := repository.Save(t.Context(), newTestTransaction())
	require.NoError(t, err)

	assert.Equal(t, 1, first.TransactionID)
	assert.Equal(t, 2, second.TransactionID)
}

func TestMemoryRepositorySaveIgnoresProvidedID(t *testing.T) {
	repository := NewMemoryRepository()
	transaction := newTestTransaction()
	transaction.TransactionID = 99

	saved, err := repository.Save(t.Context(), transaction)
	require.NoError(t, err)

	assert.Equal(t, 1, saved.TransactionID)
}

func TestMemoryRepositorySavePreservesFields(t *testing.T) {
	repository := NewMemoryRepository()
	transaction := newTestTransaction()

	saved, err := repository.Save(t.Context(), transaction)
	require.NoError(t, err)

	// The repository has no FindByID, so the stored value is checked directly.
	stored, ok := repository.transactions[saved.TransactionID]
	require.True(t, ok, "transaction %d was not stored", saved.TransactionID)

	for _, tr := range []Transaction{saved, stored} {
		assert.Equal(t, transaction.AccountID, tr.AccountID)
		assert.Equal(t, transaction.OperationTypeID, tr.OperationTypeID)
		assert.True(t, transaction.Amount.Equal(tr.Amount), "expected %s, got %s", transaction.Amount, tr.Amount)
		assert.Equal(t, transaction.EventDate, tr.EventDate)
	}
}

func TestMemoryRepositoryConcurrentSave(t *testing.T) {
	const total = 100
	repository := NewMemoryRepository()

	ids := make([]int, total)
	var wg sync.WaitGroup
	for i := range total {
		wg.Go(func() {
			tr, err := repository.Save(t.Context(), newTestTransaction())
			assert.NoError(t, err)
			ids[i] = tr.TransactionID
		})
	}
	wg.Wait()

	unique := make(map[int]bool, total)
	for _, id := range ids {
		assert.False(t, unique[id], "duplicated id %d", id)
		unique[id] = true
	}
	assert.Len(t, unique, total)
	assert.Len(t, repository.transactions, total)
}
