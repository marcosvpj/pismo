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
	if err != nil {
		return Transaction{}, ErrInvalidAccount
	}
	transaction.EventDate = s.now()
	return s.repository.Save(ctx, transaction)
}
