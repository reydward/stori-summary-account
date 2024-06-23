package handler

/*

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"stori-summary-account/summary/summary/internal/model"
	mocks "stori-summary-account/summary/summary/internal/repository/mocks"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)
func TestSummaryHandler_Summary(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSummaryRepository(ctrl)

	handler := NewSummaryHandler(mockRepo)

	tests := []struct {
		name           string
		payload        model.RequestPayload
		setupMock      func()
		expectedStatus int
		expectedBody   *model.Summary
	}{
		{
			name: "Successful summary",
			payload: model.RequestPayload{
				UserID:    1,
				AccountID: 2,
			},
			setupMock: func() {
				mockRepo.EXPECT().GetUser(gomock.Eq(1)).Return(&model.User{ID: 1, Name: "John"}, nil)
				mockRepo.EXPECT().GetAccountTotalBalance(gomock.Eq(2)).Return(float32(1000.0), nil)
				mockRepo.EXPECT().GetAccountNumberOfTransactions(gomock.Eq(2)).Return([]*model.NumberOfTransactions{{Month: "Jan", TransactionCount: 5}}, nil)
				mockRepo.EXPECT().GetAccountAverageDebit(gomock.Eq(2)).Return(float32(-100.0), nil)
				mockRepo.EXPECT().GetAccountAverageCredit(gomock.Eq(2)).Return(float32(200.0), nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: &model.Summary{
				User:                 &model.User{ID: 1, Name: "John"},
				TotalBalance:         1000.0,
				NumberOfTransactions: []*model.NumberOfTransactions{{Month: "Jan", TransactionCount: 5}},
				AverageDebitAmount:   -100.0,
				AverageCreditAmount:  200.0,
				StatusMessage:        "Email sent successfully",
			},
		},
		{
			name: "Error getting user",
			payload: model.RequestPayload{
				UserID:    1,
				AccountID: 2,
			},
			setupMock: func() {
				mockRepo.EXPECT().GetUser(gomock.Eq(1)).Return(nil, errors.New("user not found"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			payloadBytes, _ := json.Marshal(tt.payload)
			req, _ := http.NewRequest("POST", "/summary", bytes.NewBuffer(payloadBytes))
			rr := httptest.NewRecorder()

			handler.Summary(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedBody != nil {
				var gotBody model.Summary
				err := json.Unmarshal(rr.Body.Bytes(), &gotBody)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody, &gotBody)
			}
		})
	}
}
*/
