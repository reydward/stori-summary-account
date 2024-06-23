package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"stori-summary-account/summary/summary/internal/email"
	"stori-summary-account/summary/summary/internal/model"
	"stori-summary-account/summary/summary/internal/repository"
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

	user, err := h.repo.GetUser(payload.UserID)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	if user == nil {
		http.Error(writer, "User not found", http.StatusNotFound)
		return
	}

	account, err := h.repo.GetAccountInfo(payload.AccountID)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	if account == nil {
		http.Error(writer, "Account not found", http.StatusNotFound)
		return
	}

	totalBalance, err := h.repo.GetAccountTotalBalance(payload.AccountID)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	numberOfTransactions, err := h.repo.GetAccountNumberOfTransactions(payload.AccountID)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	averageDebitAmount, err := h.repo.GetAccountAverageDebit(payload.AccountID)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	averageCreditAmount, err := h.repo.GetAccountAverageCredit(payload.AccountID)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	summary := &model.Summary{
		User:                 user,
		TotalBalance:         totalBalance,
		NumberOfTransactions: numberOfTransactions,
		AverageDebitAmount:   averageDebitAmount,
		AverageCreditAmount:  averageCreditAmount,
	}

	statusMessage, err := email.SendEmail(summary)
	summary.StatusMessage = statusMessage

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(summary)
}
