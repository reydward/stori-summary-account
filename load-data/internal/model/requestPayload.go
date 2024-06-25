package model

import (
	"bytes"
)

type RequestPayload struct {
	AccountID int          `json:"accountId"`
	FileName  string       `json:"fileName"`
	File      bytes.Buffer `json:"file"`
}
