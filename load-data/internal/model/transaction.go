package model

type Transaction struct {
	ID        int     `db:"id" json:"id"`
	AccountID int     `db:"account_id" json:"accountId"`
	Date      string  `db:"date" json:"date"`
	Amount    float64 `db:"amount" json:"amount"`
}
