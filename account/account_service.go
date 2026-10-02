package account

import "errors"

var ErrNotFound = errors.New("Account not found")

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) GetAccount(accountID int) (Account, error) {
	acc, err := s.repository.GetAccount(accountID)
	if acc.AccountID == 0 && err == nil {
		return Account{}, ErrNotFound
	}
	return acc, err
}
