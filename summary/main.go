package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"log"
	"net/http"
	"stori-summary-account/summary/summary/internal/database"
	"stori-summary-account/summary/summary/internal/handler"
	"stori-summary-account/summary/summary/internal/repository"
)

func LambdaHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Starting lambda summary - ", "method: ", request.HTTPMethod, " path: ", request.Path)
	if request.HTTPMethod == "GET" && request.Path == "/summary" {
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       "Hello, World!",
		}, nil
	}
	return events.APIGatewayProxyResponse{
		StatusCode: 404,
		Body:       "Not Found",
	}, nil
}

/*
	func main() {
		lambda.Start(LambdaHandler)
	}
*/
func main() {
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
