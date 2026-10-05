package account

import "context"

type MemoryRepository struct {
	accounts []*Account
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		accounts: make([]*Account, 0),
	}
}

func (r *MemoryRepository) FindByID(ctx context.Context, accountID int) (Account, error) {
	for _, account := range r.accounts {
		if account.AccountID == accountID {
			return *account, nil
		}
	}
	return Account{}, ErrNotFound
}

func (r *MemoryRepository) Save(ctx context.Context, account Account) (Account, error) {
	account.AccountID = len(r.accounts) + 1
	r.accounts = append(r.accounts, &account)
	return account, nil
}
