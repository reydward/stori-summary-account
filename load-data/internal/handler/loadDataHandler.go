package handler

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"load-data/internal/model"
	"load-data/internal/repository"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

type LoadDataHandler struct {
	repo repository.LoadDataRepository
}

func NewLoadDataHandler(repo repository.LoadDataRepository) *LoadDataHandler {
	return &LoadDataHandler{repo: repo}
}

func (h *LoadDataHandler) Health(w http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(w, "<h1>The Lambda Load Data is working!<h1>\n")
}

func (h *LoadDataHandler) LoadData(writer http.ResponseWriter, request *http.Request) {
	// Getting the payload
	var payload model.RequestPayload
	err := json.NewDecoder(request.Body).Decode(&payload)
	if err != nil {
		http.Error(writer, "Bad Request", http.StatusBadRequest)
		return
	}

	//Open the CSV file
	file, err := os.Open(payload.FilePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	//Read the CSV file
	reader := csv.NewReader(file)
	transactionsCSV, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	//Write every record in the PostgreSQL database table
	var wg sync.WaitGroup
	resultChannel := make(chan error)
	accountID := payload.AccountID

	for _, transactionCSV := range transactionsCSV[1:] {
		transactionID, _ := strconv.Atoi(transactionCSV[0])
		date := transactionCSV[1]
		amount, _ := strconv.ParseFloat(transactionCSV[2], 64)
		transaction := model.Transaction{ID: transactionID, AccountID: accountID, Date: date, Amount: amount}

		wg.Add(1)
		go processTransaction(h.repo, transaction, &wg, resultChannel)
	}

	go func() {
		wg.Wait()
		close(resultChannel)
	}()

	responseMessage := "Data loaded successfully"
	for result := range resultChannel {
		if result != nil {
			log.Printf("Error: %v", result)
			responseMessage = "Failed to load some data, please check the logs for more information"
		}
	}

	//Setting the response
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(struct{ Message string }{responseMessage})
}

func processTransaction(repo repository.LoadDataRepository, transaction model.Transaction, wg *sync.WaitGroup, resultChannel chan<- error) {
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
