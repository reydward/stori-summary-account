package handler

import (
	"encoding/csv"
	"encoding/json"
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
	//Getting the payload
	var payload model.RequestPayload
	err := json.NewDecoder(request.Body).Decode(&payload)
	if err != nil {
		http.Error(writer, "Bad Request", http.StatusBadRequest)
		return
	}

	//Load the csv file
	file, err := os.Open(payload.FilePath)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	//Read the csv file
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Failed to read CSV file: %v", err)
	}

	//Insert every csv file row in the db table
	year := time.Now().Year()
	var wg sync.WaitGroup
	transactionChannel := make(chan model.Transaction)
	go func() {
		for transaction := range transactionChannel {
			fmt.Printf("ID: %d, Date: %s, Amount: %.2f\n", transaction.ID, transaction.Date.Format("2006-01-02"), transaction.Amount)
		}
	}()
	for _, record := range records[1:] { //without th header
		wg.Add(1)
		go func(record []string) {
			defer wg.Done()
			transaction, err := processTransaction(record, year)
			if err != nil {
				log.Printf("Failed to process record: %v", err)
				return
			}
			transactionChannel <- transaction
		}(record)
	}

	wg.Wait()
	close(transactionChannel)

	//Setting the response
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode("")
}

func processTransaction(record []string, year int) (model.Transaction, error) {
	id, err := strconv.Atoi(record[0])
	if err != nil {
		return model.Transaction{}, fmt.Errorf("invalid ID: %v", err)
	}

	date, err := time.Parse("1/2", record[1])
	if err != nil {
		return model.Transaction{}, fmt.Errorf("invalid date: %v", err)
	}

	date = date.AddDate(year, 0, 0)
	amount, err := strconv.ParseFloat(record[2], 64)
	if err != nil {
		return model.Transaction{}, fmt.Errorf("invalid amount: %v", err)
	}

	return model.Transaction{
		ID:     id,
		Date:   date,
		Amount: amount,
	}, nil
}
