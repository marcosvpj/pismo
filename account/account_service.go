package account

import "errors"

type Service struct {
	repository Repository
}

var ErrFieldDocumentNumberMissing = errors.New("field document_number is required")

func NewService(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) GetAccount(accountID int) (Account, error) {
	return s.repository.FindByID(accountID)
}

func (s *Service) CreateAccount(account Account) (Account, error) {
	if account.DocumentNumber == "" {
		return Account{}, ErrFieldDocumentNumberMissing
	}

	return s.repository.SaveAccount(account)
}
