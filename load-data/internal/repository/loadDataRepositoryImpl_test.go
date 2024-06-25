package repository

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"load-data/internal/model"
	"testing"
)

func TestInsertTransaction(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewLoadDataRepository(db)

	transaction := model.Transaction{
		ID:        0,
		AccountID: 1,
		Date:      "2024-07-15",
		Amount:    60.5,
	}

	mock.ExpectQuery(`INSERT INTO transactions \(id, account_id, date, amount\) VALUES \(\$1, \$2, \$3, \$4\) RETURNING id`).
		WithArgs(transaction.ID, transaction.AccountID, transaction.Date, transaction.Amount).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	result, err := repo.InsertTransaction(transaction)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.ID)
}

func TestInsertTransactionError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewLoadDataRepository(db)

	transaction := model.Transaction{
		ID:        0,
		AccountID: 1,
		Date:      "2024-07-15",
		Amount:    60.5,
	}

	mock.ExpectQuery(`INSERT INTO transactions \(id, account_id, date, amount\) VALUES \(\$1, \$2, \$3, \$4\) RETURNING id`).
		WithArgs(transaction.ID, transaction.AccountID, transaction.Date, transaction.Amount).
		WillReturnError(sql.ErrConnDone)

	result, err := repo.InsertTransaction(transaction)

	assert.Error(t, err)
	assert.Nil(t, result)
}
