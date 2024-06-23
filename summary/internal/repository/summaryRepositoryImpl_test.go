package repository

import (
	"database/sql"
	"errors"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewSummaryRepository(db)

	query := "SELECT \\* FROM users WHERE id = \\$1"
	rows := sqlmock.NewRows([]string{"id", "name", "email"}).AddRow(1, "John Doe", "john@example.com")

	mock.ExpectQuery(query).WithArgs(1).WillReturnRows(rows)

	user, err := repo.GetUser(1)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, 1, user.ID)
	assert.Equal(t, "John Doe", user.Name)
	assert.Equal(t, "john@example.com", user.Email)
}

func TestGetUser_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery("SELECT \\* FROM users WHERE id = \\$1").
		WithArgs(123).
		WillReturnError(sql.ErrNoRows)

	repo := &summaryRepositoryImpl{db: db}
	user, err := repo.GetUser(123)

	assert.Nil(t, user)
	assert.Nil(t, err)
}

func TestGetAccountInfo(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewSummaryRepository(db)

	query := "SELECT \\* FROM accounts WHERE id = \\$1"
	rows := sqlmock.NewRows([]string{"id", "user_id", "name"}).AddRow(1, 1, "Main Account")

	mock.ExpectQuery(query).WithArgs(1).WillReturnRows(rows)

	account, err := repo.GetAccountInfo(1)
	assert.NoError(t, err)
	assert.NotNil(t, account)
	assert.Equal(t, 1, account.ID)
	assert.Equal(t, 1, account.UserID)
	assert.Equal(t, "Main Account", account.Name)
}

func TestGetAccountInfo_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery("SELECT \\* FROM accounts WHERE id = \\$1").
		WithArgs(123).
		WillReturnError(sql.ErrNoRows)

	repo := &summaryRepositoryImpl{db: db}
	user, err := repo.GetAccountInfo(123)

	assert.Nil(t, user)
	assert.Nil(t, err)
}

func TestGetAccountAverageDebit(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewSummaryRepository(db)

	query := `
        SELECT AVG\(amount\) AS average_debit_amount
        FROM transactions
        WHERE account_id = \$1 AND amount < 0;
    `
	rows := sqlmock.NewRows([]string{"average_debit_amount"}).AddRow(-20.0)

	mock.ExpectQuery(query).WithArgs(1).WillReturnRows(rows)

	avgDebit, err := repo.GetAccountAverageDebit(1)
	assert.NoError(t, err)
	assert.Equal(t, float32(-20.0), avgDebit)
}

func TestGetAccountAverageDebit_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery("SELECT AVG\\(amount\\) AS average_debit_amount FROM transactions WHERE account_id = \\$1 AND amount < 0;").
		WithArgs(111).
		WillReturnError(errors.New("database error"))

	repo := &summaryRepositoryImpl{db: db}

	avg, err := repo.GetAccountAverageDebit(111)
	assert.Error(t, err)
	assert.Equal(t, "error getting the average debit amount: database error", err.Error())
	assert.Equal(t, float32(0), avg)
}

func TestGetAccountAverageCredit(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewSummaryRepository(db)

	query := `
        SELECT AVG\(amount\) AS average_credit_amount
        FROM transactions
        WHERE account_id = \$1 AND amount > 0;
    `
	rows := sqlmock.NewRows([]string{"average_credit_amount"}).AddRow(30.0)

	mock.ExpectQuery(query).WithArgs(1).WillReturnRows(rows)

	avgCredit, err := repo.GetAccountAverageCredit(1)
	assert.NoError(t, err)
	assert.Equal(t, float32(30.0), avgCredit)
}

func TestGetAccountAverageCredit_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery("SELECT AVG\\(amount\\) AS average_credit_amount FROM transactions WHERE account_id = \\$1 AND amount > 0;").
		WithArgs(111).
		WillReturnError(errors.New("database error"))

	repo := &summaryRepositoryImpl{db: db}

	avg, err := repo.GetAccountAverageCredit(111)
	assert.Error(t, err)
	assert.Equal(t, "error getting the average credit amount: database error", err.Error())
	assert.Equal(t, float32(0), avg)
}

func TestGetAccountNumberOfTransactions(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewSummaryRepository(db)

	query := `
        SELECT TO_CHAR\(DATE_TRUNC\('month', date\), 'Month'\) AS month, COUNT\(\*\) AS transaction_count
        FROM transactions
        WHERE account_id = \$1
        GROUP BY month
        ORDER BY month;
    `
	rows := sqlmock.NewRows([]string{"month", "transaction_count"}).AddRow("January", 5).AddRow("February", 3)

	mock.ExpectQuery(query).WithArgs(1).WillReturnRows(rows)

	transactions, err := repo.GetAccountNumberOfTransactions(1)
	assert.NoError(t, err)
	assert.Len(t, transactions, 2)
	assert.Equal(t, "January", transactions[0].Month)
	assert.Equal(t, 5, transactions[0].TransactionCount)
	assert.Equal(t, "February", transactions[1].Month)
	assert.Equal(t, 3, transactions[1].TransactionCount)
}

func TestGetAccountNumberOfTransactions_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery("SELECT TO_CHAR.*").WillReturnError(errors.New("db error"))
	repo := &summaryRepositoryImpl{db: db}

	accountID := 111
	result, err := repo.GetAccountNumberOfTransactions(accountID)

	assert.Error(t, err)
	assert.Nil(t, result)

	expectedError := "error executing the query: db error"
	assert.EqualError(t, err, expectedError)
}

func TestGetAccountTotalBalance(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewSummaryRepository(db)

	query := "SELECT COALESCE\\(SUM\\(amount\\), 0\\) AS total_balance FROM transactions WHERE account_id = \\$1"
	rows := sqlmock.NewRows([]string{"total_balance"}).AddRow(100.0)

	mock.ExpectQuery(query).WithArgs(1).WillReturnRows(rows)

	totalBalance, err := repo.GetAccountTotalBalance(1)
	assert.NoError(t, err)
	assert.Equal(t, float32(100.0), totalBalance)
}

func TestGetAccountTotalBalance_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery("SELECT COALESCE\\(SUM\\(amount\\), 0\\) AS total_balance FROM transactions WHERE account_id = \\$1").
		WithArgs(111).
		WillReturnError(errors.New("db error"))
	repo := &summaryRepositoryImpl{db: db}
	totalBalance, err := repo.GetAccountTotalBalance(111)

	assert.Error(t, err)

	expectedError := "error getting the total balance: db error"
	assert.EqualError(t, err, expectedError)
	assert.Equal(t, float32(0), totalBalance)
}
