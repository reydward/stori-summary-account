package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"summary/internal/email"
	"summary/internal/model"
	"summary/internal/repository"
)

type SummaryHandler struct {
	repo repository.SummaryRepository
}

func NewSummaryHandler(repo repository.SummaryRepository) *SummaryHandler {
	return &SummaryHandler{repo: repo}
}

func (h *SummaryHandler) Health(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "<h1>The Lambda Summary is working!<h1>\n")
}

func (h *SummaryHandler) Summary(writer http.ResponseWriter, request *http.Request) {
	var payload model.RequestPayload
	err := json.NewDecoder(request.Body).Decode(&payload)
	if err != nil {
		http.Error(writer, "Bad Request", http.StatusBadRequest)
		return
	}
	//Getting user information
	user, err := h.repo.GetUser(payload.UserID)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	if user == nil {
		http.Error(writer, "User not found", http.StatusNotFound)
		return
	}
	//Getting account information
	account, err := h.repo.GetAccountInfo(payload.AccountID)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	if account == nil {
		http.Error(writer, "Account not found", http.StatusNotFound)
		return
	}
	//Getting total balance
	totalBalance, err := h.repo.GetAccountTotalBalance(payload.AccountID)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	//Getting number of transactions per month
	numberOfTransactions, err := h.repo.GetAccountNumberOfTransactions(payload.AccountID)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	//Getting average debit amount
	averageDebitAmount, err := h.repo.GetAccountAverageDebit(payload.AccountID)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	//Getting average credit amount
	averageCreditAmount, err := h.repo.GetAccountAverageCredit(payload.AccountID)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	//Building the summary object
	summary := &model.Summary{
		User:                 user,
		TotalBalance:         totalBalance,
		NumberOfTransactions: numberOfTransactions,
		AverageDebitAmount:   averageDebitAmount,
		AverageCreditAmount:  averageCreditAmount,
	}
	//Sending the email
	statusMessage, err := email.SendEmail(summary)
	summary.StatusMessage = statusMessage
	//Setting the response
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(summary)
}
