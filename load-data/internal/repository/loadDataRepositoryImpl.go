package repository

import (
	"database/sql"
	"load-data/internal/model"
)

type loadDataRepositoryImpl struct {
	db *sql.DB
}

func (l *loadDataRepositoryImpl) InsertTransaction(transaction model.Transaction) (*model.Transaction, error) {
	query := `
        INSERT INTO transactions (id, account_id, date, amount)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `
	err := l.db.QueryRow(query, transaction.ID, transaction.AccountID, transaction.Date, transaction.Amount).Scan(&transaction.ID)
	if err != nil {
		return nil, err
	}

	return &transaction, nil
}

func NewLoadDataRepository(db *sql.DB) LoadDataRepository {
	return &loadDataRepositoryImpl{db: db}
}
