package service

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"load-data/internal/file"
	"load-data/internal/model"
	"load-data/internal/repository"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

func ProcessTransactions(payload model.RequestPayload, repository repository.LoadDataRepository) (string, error) {
	log.Printf("Processing the transactions for account ID: %d", payload.AccountID)

	//Setting the AWS configuration
	err, awsService := setAWSConfiguration()
	if err != nil {
		log.Fatalf("error setting the AWS configuration: %v", err)
	}

	//Upload file to S3
	err = awsService.UploadFile(os.Getenv("S3_BUCKET_NAME"), payload.FileName, payload.File)
	if err != nil {
		log.Fatalf("error uploading the file to S3: %v", err)
	}
	log.Printf("File uploaded successfully to S3: %s", payload.FileName)

	//Get the CSV file from S3
	s3File := awsService.GetFile(os.Getenv("S3_BUCKET_NAME"), payload.FileName)
	log.Printf("File got successfully from S3: %s", s3File)

	//Read the CSV file
	reader := csv.NewReader(s3File)
	transactionsCSV, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	//Write every record in the PostgreSQL database table
	resultChannel := writeTransaction(payload, repository, transactionsCSV)

	//Build the response message
	responseMessage, errorResult := buildResult(resultChannel)

	//Delete the CSV file from S3
	err = awsService.DeleteFile(os.Getenv("S3_BUCKET_NAME"), payload.FileName)
	if err != nil {
		log.Fatalf("error deleting the file from S3: %v", err)
	}
	return responseMessage, errorResult
}

func setAWSConfiguration() (error, file.AWSService) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(os.Getenv("S3_REGION")))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	awsService := file.AWSService{
		S3Client: s3.NewFromConfig(cfg),
	}
	return err, awsService
}

func buildResult(resultChannel chan error) (string, error) {
	responseMessage := "Data loaded successfully"
	var errorResult error
	errorResult = nil

	for result := range resultChannel {
		if result != nil {
			log.Printf("Error: %v", result)
			responseMessage = "Failed to load some data, please check the logs for more information"
			errorResult = errors.New(responseMessage)
		}
	}
	return responseMessage, errorResult
}

func writeTransaction(payload model.RequestPayload, repository repository.LoadDataRepository, transactionsCSV [][]string) chan error {
	var wg sync.WaitGroup
	resultChannel := make(chan error)
	accountID := payload.AccountID

	for _, transactionCSV := range transactionsCSV[1:] {
		transactionID, _ := strconv.Atoi(transactionCSV[0])
		date := transactionCSV[1]
		amount, _ := strconv.ParseFloat(transactionCSV[2], 64)
		transaction := model.Transaction{ID: transactionID, AccountID: accountID, Date: date, Amount: amount}

		wg.Add(1)
		go insertTransaction(repository, transaction, &wg, resultChannel)
	}

	go func() {
		wg.Wait()
		close(resultChannel)
	}()
	return resultChannel
}

func insertTransaction(repo repository.LoadDataRepository, transaction model.Transaction, wg *sync.WaitGroup, resultChannel chan<- error) {
	defer wg.Done()

	date, err := time.Parse("1/2", transaction.Date)
	if err != nil {
		log.Printf("Failed to parse the date: %v", err)
	}

	year := time.Now().Year()
	date = date.AddDate(year, 0, 0)
	transaction.Date = date.Format("2006-01-02")

	fmt.Printf("ID: %d, AccountID: %d, Date: %s, Transaction: %.2f\n", transaction.ID, transaction.AccountID, transaction.Date, transaction.Amount)

	//Insert the transaction in the database
	_, err = repo.InsertTransaction(transaction)
	if err != nil {
		resultChannel <- errors.New("Failed to insert the transaction with ID: " + strconv.Itoa(transaction.ID) + " error: " + err.Error())
	}

	resultChannel <- nil
}
