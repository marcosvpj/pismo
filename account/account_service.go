package account

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) GetAccount(accountID int) (Account, error) {
	return s.repository.GetAccount(accountID)
}
