package account

type Repository interface {
	GetAccount(accountID int) (Account, error)
}
