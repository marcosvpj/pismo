package account

import (
	"context"
	"sync"
)

type MemoryRepository struct {
	mu       sync.RWMutex
	accounts map[int]Account
	nextID   int
}

var _ Repository = (*MemoryRepository)(nil)

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		accounts: make(map[int]Account),
		nextID:   1,
	}
}

func (r *MemoryRepository) FindByID(ctx context.Context, accountID int) (Account, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	acc, ok := r.accounts[accountID]
	if !ok {
		return Account{}, ErrNotFound
	}

	return acc, nil
}

func (r *MemoryRepository) Save(ctx context.Context, account Account) (Account, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	account.AccountID = r.nextID
	r.accounts[r.nextID] = account
	r.nextID++
	return account, nil
}
