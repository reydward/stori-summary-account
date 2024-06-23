package repository

import "stori-summary-account/summary/summary/internal/model"

type SummaryRepository interface {
	GetUser(userID int) (*model.User, error)
	GetAccountInfo(accountID int) (*model.Account, error)
	GetAccountTotalBalance(accountID int) (float32, error)
	GetAccountNumberOfTransactions(accountID int) ([]*model.NumberOfTransactions, error)
	GetAccountAverageDebit(accountID int) (float32, error)
	GetAccountAverageCredit(accountID int) (float32, error)
}
