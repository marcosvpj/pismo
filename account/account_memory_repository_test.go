package account

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemoryRepositorySaveAssignsSequentialIDs(t *testing.T) {
	repository := NewMemoryRepository()

	first, err := repository.Save(t.Context(), Account{DocumentNumber: "11111111111"})
	require.NoError(t, err)
	second, err := repository.Save(t.Context(), Account{DocumentNumber: "22222222222"})
	require.NoError(t, err)

	assert.Equal(t, Account{AccountID: 1, DocumentNumber: "11111111111"}, first)
	assert.Equal(t, Account{AccountID: 2, DocumentNumber: "22222222222"}, second)
}

func TestMemoryRepositorySaveIgnoresProvidedID(t *testing.T) {
	repository := NewMemoryRepository()

	saved, err := repository.Save(t.Context(), Account{AccountID: 99, DocumentNumber: "12345678900"})
	require.NoError(t, err)

	assert.Equal(t, 1, saved.AccountID)
}

func TestMemoryRepositoryFindByID(t *testing.T) {
	repository := NewMemoryRepository()
	first, _ := repository.Save(t.Context(), Account{DocumentNumber: "11111111111"})
	second, _ := repository.Save(t.Context(), Account{DocumentNumber: "22222222222"})

	tests := []struct {
		name        string
		accountID   int
		expected    Account
		expectedErr error
	}{
		{name: "First account", accountID: 1, expected: first},
		{name: "Second account", accountID: 2, expected: second},
		{name: "Not found", accountID: 3, expected: Account{}, expectedErr: ErrNotFound},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			acc, err := repository.FindByID(t.Context(), test.accountID)

			assert.ErrorIs(t, err, test.expectedErr)
			assert.Equal(t, test.expected, acc)
		})
	}
}

func TestMemoryRepositoryConcurrentSave(t *testing.T) {
	const total = 100
	repository := NewMemoryRepository()

	ids := make([]int, total)
	var wg sync.WaitGroup
	for i := range total {
		wg.Go(func() {
			acc, err := repository.Save(t.Context(), Account{DocumentNumber: "12345678900"})
			assert.NoError(t, err)
			ids[i] = acc.AccountID
		})
	}
	wg.Wait()

	unique := make(map[int]bool, total)
	for _, id := range ids {
		assert.False(t, unique[id], "duplicated id %d", id)
		unique[id] = true

		_, err := repository.FindByID(t.Context(), id)
		assert.NoError(t, err)
	}
	assert.Len(t, unique, total)
}
