package model

type RequestPayload struct {
	AccountID int    `json:"accountId"`
	FilePath  string `json:"filePath"`
}
