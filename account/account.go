package account

type Account struct {
	AccountID      int    `json:"account_id,omitempty"`
	DocumentNumber string `json:"document_number"`
}
