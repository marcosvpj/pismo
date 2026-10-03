package account

import (
	"context"
	"errors"
)

type Service struct {
	repository Repository
}

var ErrFieldDocumentNumberMissing = errors.New("field document_number is required")

func NewService(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) GetAccount(ctx context.Context, accountID int) (Account, error) {
	return s.repository.FindByID(ctx, accountID)
}

func (s *Service) CreateAccount(ctx context.Context, account Account) (Account, error) {
	if account.DocumentNumber == "" {
		return Account{}, ErrFieldDocumentNumberMissing
	}

	return s.repository.Save(ctx, account)
}
