package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"load-data/internal/database"
	"load-data/internal/handler"
	"load-data/internal/model"
	"load-data/internal/repository"
	"log"
	"net/http"
)

var db *sql.DB

func init() {
	db, err := database.NewPostgresConnection()
	if err != nil {
		fmt.Printf("Failed to connect to the database: %v\n", err)
		return
	}
	defer db.Close()
}

func lambdaHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Starting lambda summary function")

	if request.HTTPMethod == "POST" && request.Path == "/load-data" {

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
	}

	var apigwresponse = &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "",
	}

	apigwresponse.Headers = make(map[string]string)
	apigwresponse.Headers["Access-Control-Allow-Origin"] = "*"
	apigwresponse.Headers["Access-Control-Allow-Methods"] = "GET,POST,OPTIONS"
	return *apigwresponse, nil
}

func main() {
	//	lambda.Start(lambdaHandler)

	summaryRepository := repository.NewSummaryRepository(db)
	summaryHandler := handler.NewSummaryHandler(summaryRepository)

	http.HandleFunc("/", summaryHandler.Health)
	http.HandleFunc("/summary", summaryHandler.Summary)
	http.ListenAndServe(":3000", nil)

}
