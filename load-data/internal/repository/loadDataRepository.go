package repository

import "load-data/internal/model"

type LoadDataRepository interface {
	InsertTransaction(transaction model.Transaction) (*model.Transaction, error)
}
