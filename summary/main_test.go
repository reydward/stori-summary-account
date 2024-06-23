package main

import (
	"errors"
	"summary/internal/repository/mocks"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"summary/internal/model"
)

func TestGetSummaryData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSummaryRepository(ctrl)

	// Define los valores de retorno esperados
	user := &model.User{ID: 1, Name: "John Doe", Email: "john.doe@example.com"}
	account := &model.Account{ID: 1, UserID: 1, Name: "John's Account"}
	totalBalance := float32(1000.0)
	numberOfTransactions := []*model.NumberOfTransactions{
		{Month: "January", TransactionCount: 5},
	}
	averageDebitAmount := float32(-50.0)
	averageCreditAmount := float32(100.0)

	payload := model.RequestPayload{UserID: 1, AccountID: 1}

	mockRepo.EXPECT().GetUser(payload.UserID).Return(user, nil)
	mockRepo.EXPECT().GetAccountInfo(payload.AccountID).Return(account, nil)
	mockRepo.EXPECT().GetAccountTotalBalance(payload.AccountID).Return(totalBalance, nil)
	mockRepo.EXPECT().GetAccountNumberOfTransactions(payload.AccountID).Return(numberOfTransactions, nil)
	mockRepo.EXPECT().GetAccountAverageDebit(payload.AccountID).Return(averageDebitAmount, nil)
	mockRepo.EXPECT().GetAccountAverageCredit(payload.AccountID).Return(averageCreditAmount, nil)

	summary, err := getSummaryData(mockRepo, payload)

	assert.NoError(t, err)
	assert.Equal(t, user, summary.User)
	assert.Equal(t, totalBalance, summary.TotalBalance)
	assert.Equal(t, numberOfTransactions, summary.NumberOfTransactions)
	assert.Equal(t, averageDebitAmount, summary.AverageDebitAmount)
	assert.Equal(t, averageCreditAmount, summary.AverageCreditAmount)
}

func TestGetSummaryData_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSummaryRepository(ctrl)

	payload := model.RequestPayload{UserID: 1, AccountID: 1}

	// Configurar el comportamiento de los métodos simulados
	mockRepo.EXPECT().GetUser(payload.UserID).Return(nil, nil)

	// Llamar al método que se está probando
	summary, err := getSummaryData(mockRepo, payload)

	// Validar los resultados
	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())
	assert.Equal(t, model.Summary{}, summary)
}

func TestGetSummaryData_AccountNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSummaryRepository(ctrl)

	payload := model.RequestPayload{UserID: 1, AccountID: 1}

	// Configurar el comportamiento de los métodos simulados
	mockRepo.EXPECT().GetUser(payload.UserID).Return(&model.User{ID: 1, Name: "John Doe", Email: "john.doe@example.com"}, nil)
	mockRepo.EXPECT().GetAccountInfo(payload.AccountID).Return(nil, nil)

	// Llamar al método que se está probando
	summary, err := getSummaryData(mockRepo, payload)

	// Validar los resultados
	assert.Error(t, err)
	assert.Equal(t, "account not found", err.Error())
	assert.Equal(t, model.Summary{}, summary)
}

func TestGetSummaryData_OtherErrors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSummaryRepository(ctrl)

	payload := model.RequestPayload{UserID: 1, AccountID: 1}

	gomock.InOrder(
		mockRepo.EXPECT().GetUser(payload.UserID).Return(&model.User{ID: 1, Name: "John Doe", Email: "john.doe@example.com"}, nil),
		mockRepo.EXPECT().GetAccountInfo(payload.AccountID).Return(&model.Account{ID: 1, UserID: 1, Name: "John's Account"}, nil),
		mockRepo.EXPECT().GetAccountTotalBalance(payload.AccountID).Return(float32(0), errors.New("some error")),
	)
	mockRepo.EXPECT().GetAccountNumberOfTransactions(payload.AccountID).Times(1)
	mockRepo.EXPECT().GetAccountAverageDebit(payload.AccountID).Times(1)
	mockRepo.EXPECT().GetAccountAverageCredit(payload.AccountID).Times(1)

	// Llamar al método que se está probando
	summary, err := getSummaryData(mockRepo, payload)

	// Validar los resultados
	assert.Error(t, err)
	assert.Equal(t, "error getting the summary data", err.Error())
	assert.Equal(t, model.Summary{}, summary)
}
