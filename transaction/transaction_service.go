package transaction

import (
	"context"
	"errors"
	"time"

	"github.com/marcosvpj/pismo/account"
)

type Service struct {
	repository        Repository
	accountRepository account.Repository
	now               func() time.Time
}

var ErrInvalidAccount = errors.New("invalid account")

func NewService(repository Repository, accountRepository account.Repository, now func() time.Time) *Service {
	return &Service{
		repository:        repository,
		accountRepository: accountRepository,
		now:               now,
	}
}

func (s *Service) CreateTransaction(ctx context.Context, transaction Transaction) (Transaction, error) {
	_, err := s.accountRepository.FindByID(ctx, transaction.AccountID)
	if errors.Is(err, account.ErrNotFound) {
		return Transaction{}, ErrInvalidAccount
	} else if err != nil {
		return Transaction{}, err
	}

	transaction.EventDate = s.now()
	return s.repository.Save(ctx, transaction)
}
