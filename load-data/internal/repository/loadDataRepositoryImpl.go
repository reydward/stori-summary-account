package repository

import (
	"database/sql"
)

type loadDataRepositoryImpl struct {
	db *sql.DB
}

func NewLoadDataRepository(db *sql.DB) LoadDataRepository {
	return &loadDataRepositoryImpl{db: db}
}
