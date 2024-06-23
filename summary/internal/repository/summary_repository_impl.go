package repository

import (
	"database/sql"
	"fmt"
	"stori-summary-account/summary/summary/internal/model"
)

type summaryRepositoryImpl struct {
	db *sql.DB
}

func (r *summaryRepositoryImpl) GetUser(userID int) (*model.User, error) {
	return &model.User{
		ID:    userID,
		Name:  "John Doe",
		Email: "info@brainsmartsolutions.com",
	}, nil
}

func (r *summaryRepositoryImpl) GetAccountAverageDebit(accountID int) (float32, error) {
	var averageDebitAmount sql.NullFloat64 //in order to handle null values

	query := `
        SELECT AVG(amount) AS average_debit_amount
        FROM transactions
        WHERE account_id = $1 AND amount < 0;
    `

	err := r.db.QueryRow(query, accountID).Scan(&averageDebitAmount)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, fmt.Errorf("error getting the average debit amount: %v", err)
	}

	if averageDebitAmount.Valid {
		return float32(averageDebitAmount.Float64), nil
	}

	return 0, nil
}

func (r *summaryRepositoryImpl) GetAccountAverageCredit(accountID int) (float32, error) {
	var averageDebitAmount sql.NullFloat64 //in order to handle null values

	query := `
        SELECT AVG(amount) AS average_credit_amount
        FROM transactions
        WHERE account_id = $1 AND amount > 0;
    `

	err := r.db.QueryRow(query, accountID).Scan(&averageDebitAmount)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, fmt.Errorf("error getting the average credit amount: %v", err)
	}

	if averageDebitAmount.Valid {
		return float32(averageDebitAmount.Float64), nil
	}

	return 0, nil
}

func (r *summaryRepositoryImpl) GetAccountNumberOfTransactions(accountID int) ([]*model.NumberOfTransactions, error) {
	var results []*model.NumberOfTransactions
	query := `
        SELECT TO_CHAR(DATE_TRUNC('month', date), 'Month') AS month, COUNT(*) AS transaction_count
        FROM transactions
        WHERE account_id = $1
        GROUP BY month
        ORDER BY month;
    `
	rows, err := r.db.Query(query, accountID)
	if err != nil {
		return nil, fmt.Errorf("error executing the query: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var numberTransations model.NumberOfTransactions
		err := rows.Scan(&numberTransations.Month, &numberTransations.TransactionCount)
		if err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}
		results = append(results, &numberTransations)
	}

	return results, nil
}

func (s summaryRepositoryImpl) GetAccountTotalBalance(accountID int) (float32, error) {
	var totalBalance float32

	query := "SELECT COALESCE(SUM(amount), 0) AS total_balance FROM transactions WHERE account_id = $1"
	err := s.db.QueryRow(query, accountID).Scan(&totalBalance)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, fmt.Errorf("error getting the total amount: %v", err)
	}

	return totalBalance, nil
}

func NewSummaryRepository(db *sql.DB) SummaryRepository {
	return &summaryRepositoryImpl{db: db}
}
