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
		accounts: make(map[int]Account, 0),
		nextID:   1,
	}
}

func (r *MemoryRepository) FindByID(ctx context.Context, accountID int) (Account, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, account := range r.accounts {
		if account.AccountID == accountID {
			return account, nil
		}
	}
	return Account{}, ErrNotFound
}

func (r *MemoryRepository) Save(ctx context.Context, account Account) (Account, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	account.AccountID = r.nextID
	r.nextID++
	r.accounts[r.nextID] = account
	return account, nil
}
