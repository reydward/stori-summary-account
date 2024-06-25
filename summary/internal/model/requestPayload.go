package model

type RequestPayload struct {
	UserID    int    `json:"user"`
	AccountID int    `json:"accountId"`
	Quarter   string `json:"quarter"`
}
