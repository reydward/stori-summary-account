package model

type Summary struct {
	User                 *User                   `json:"user"`
	TotalBalance         float32                 `json:"totalBalance"`
	NumberOfTransactions []*NumberOfTransactions `json:"numberOfTransactions"`
	AverageDebitAmount   float32                 `json:"averageDebitAmount"`
	AverageCreditAmount  float32                 `json:"averageCreditAmount"`
	StatusMessage        string                  `json:"statusMessage"`
}

type NumberOfTransactions struct {
	Month            string `db:"month"`
	TransactionCount int    `db:"transaction_count"`
}
