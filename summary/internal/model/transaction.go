package model

type Transaction struct {
	ID        int    `db:"id" json:"id"`
	AccountID int    `db:"account_id" json:"accountId"`
	date      string `db:"date" json:"date"`
	amount    int    `db:"amount" json:"amount"`
}
