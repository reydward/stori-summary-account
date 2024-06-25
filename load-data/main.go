package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"load-data/internal/database"
	"load-data/internal/handler"
	"load-data/internal/repository"
	"net/http"
)

func serverExecution() {
	db, err := database.NewPostgresConnection()
	if err != nil {
		fmt.Printf("Failed to connect to the database: %v\n", err)
		return
	}
	defer db.Close()
	loadDataRepository := repository.NewLoadDataRepository(db)
	loadDataHandler := handler.NewLoadDataHandler(loadDataRepository)

	http.HandleFunc("/", loadDataHandler.Health)
	http.HandleFunc("/load-data", loadDataHandler.LoadData)
	http.ListenAndServe(":4000", nil)
}

func main() {
	lambda.Start(handler.LambdaHandler)
	//serverExecution()
}
