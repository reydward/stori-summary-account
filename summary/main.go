package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"log"
	"net/http"
	"summary/internal/database"
	"summary/internal/email"
	"summary/internal/handler"
	"summary/internal/model"
	"summary/internal/repository"
)

func lambdaHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	db, err := database.NewPostgresConnection()
	if err != nil {
		log.Printf("Failed to connect to the database: %v\n", err)
	}
	defer db.Close()

	log.Printf("Starting lambda summary function")
	log.Printf("Request.HTTPMethod: %v", request.HTTPMethod)
	var response string

	if request.HTTPMethod == "POST" {
		log.Printf("lambdaHandler.db: %v", db)
		repo := repository.NewSummaryRepository(db)

		//Getting the payload
		var payload model.RequestPayload
		err := json.Unmarshal([]byte(request.Body), &payload)
		if err != nil {
			log.Printf("Failed to unmarshal request body: %v", err)
			return events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:       fmt.Sprintf(`{"error": "Invalid request payload: %v"}`, err),
			}, nil
		}

		//Getting the summary data
		summary, err := getSummaryData(repo, payload)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       fmt.Sprintf(`{"error": "Failed getting the summary data: %v"}`, err),
			}, nil
		}

		//Sending the email
		statusMessage, err := email.SendEmail(&summary)
		summary.StatusMessage = statusMessage
		log.Printf("Summary: %v", summary)

		//Setting the response
		summaryJSON, err := json.Marshal(summary)
		if err != nil {
			log.Fatalf("Error al convertir el objeto a JSON: %v", err)
		}
		response = string(summaryJSON)
	}

	log.Printf("Response: %v", response)

	var apigwresponse = &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       response,
	}
	apigwresponse.Headers = make(map[string]string)
	apigwresponse.Headers["Access-Control-Allow-Origin"] = "*"
	apigwresponse.Headers["Access-Control-Allow-Methods"] = "POST,OPTIONS"
	return *apigwresponse, nil
}

func getSummaryData(repo repository.SummaryRepository, payload model.RequestPayload) (model.Summary, error) {
	log.Printf("getSummaryData.payload: %v", payload)

	var summary = model.Summary{}
	var err error
	var isThereError bool

	//Getting user information
	user, err := repo.GetUser(payload.UserID)
	if err != nil {
		return summary, err
	}
	if user == nil {
		return summary, errors.New("user not found")
	}
	//Getting account information
	account, err := repo.GetAccountInfo(payload.AccountID)
	if err != nil {
		return summary, err
	}
	if account == nil {
		return summary, errors.New("account not found")
	}
	//Getting total balance
	totalBalance, err := repo.GetAccountTotalBalance(payload.AccountID)
	if err != nil {
		log.Printf("Failed getting the total balance: %v", err)
		isThereError = true
	}
	//Getting number of transactions per month
	numberOfTransactions, err := repo.GetAccountNumberOfTransactions(payload.AccountID)
	if err != nil {
		log.Printf("Failed getting number of transactions per month: %v", err)
		isThereError = true
	}
	//Getting average debit amount
	averageDebitAmount, err := repo.GetAccountAverageDebit(payload.AccountID)
	if err != nil {
		log.Printf("Failed getting the average debit amount: %v", err)
		isThereError = true
	}
	//Getting average credit amount
	averageCreditAmount, err := repo.GetAccountAverageCredit(payload.AccountID)
	if err != nil {
		log.Printf("Failed getting the average credit amount: %v", err)
		isThereError = true
	}

	if isThereError {
		return summary, errors.New("error getting the summary data")
	}

	//Building the summary object
	summary = model.Summary{
		User:                 user,
		Account:              account,
		TotalBalance:         totalBalance,
		NumberOfTransactions: numberOfTransactions,
		AverageDebitAmount:   averageDebitAmount,
		AverageCreditAmount:  averageCreditAmount,
	}

	log.Printf("getSummaryData.summary: %v", summary)
	return summary, nil
}

func main() {
	//lambda.Start(lambdaHandler)
	httpServerExecution()
}

func httpServerExecution() {
	db, err := database.NewPostgresConnection()
	if err != nil {
		fmt.Printf("Failed to connect to the database: %v\n", err)
		return
	}
	defer db.Close()

	summaryRepository := repository.NewSummaryRepository(db)
	summaryHandler := handler.NewSummaryHandler(summaryRepository)

	http.HandleFunc("/", summaryHandler.Health)
	http.HandleFunc("/summary", summaryHandler.Summary)
	http.ListenAndServe(":3000", nil)
}
