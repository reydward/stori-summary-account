package model

type Account struct {
	ID     int    `db:"id" json:"id"`
	UserID string `db:"user_id" json:"userId"`
	Name   string `db:"name" json:"name"`
}
