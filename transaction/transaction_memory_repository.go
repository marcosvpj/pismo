package transaction

import (
	"context"
	"sync"
)

type MemoryRepository struct {
	mu           sync.RWMutex
	transactions map[int]Transaction
	nextID       int
}

var _ Repository = (*MemoryRepository)(nil)

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		transactions: make(map[int]Transaction),
		nextID:       1,
	}
}

func (r *MemoryRepository) Save(ctx context.Context, transaction Transaction) (Transaction, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	transaction.TransactionID = r.nextID
	r.transactions[r.nextID] = transaction
	r.nextID++
	return transaction, nil
}
