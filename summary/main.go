package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"log"
	"net/http"
	"net/smtp"
	"stori-summary-account/summary/summary/internal/database"
	"stori-summary-account/summary/summary/internal/handler"
	"stori-summary-account/summary/summary/internal/repository"
	"strconv"
	"strings"
)

type RequestPayload struct {
	UserID    int    `json:"user"`
	AccountID int    `json:"accountId"`
	Quarter   string `json:"quarter"`
}

type User struct {
	ID    int
	Name  string
	Email string
}

type Transaction struct {
	ID          int
	AccountID   int
	Date        string
	Transaction float64
}

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

func health(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "<h1>The Lambda Summary is working!<h1>\n")
}

func summary(writer http.ResponseWriter, request *http.Request) {
	log.Printf("Starting lambda summary - ", "method: ", request.Method)
	if request.Method != http.MethodPost {
		http.Error(writer, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload RequestPayload
	err := json.NewDecoder(request.Body).Decode(&payload)
	if err != nil {
		http.Error(writer, "Bad Request", http.StatusBadRequest)
		return
	}
	fmt.Println("payload", payload)
	/*
		db, err := sql.Open("postgres", "user=username password=password dbname=database_name sslmode=disable")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		var user User
		err = db.QueryRow("SELECT id, name, email FROM users WHERE id = $1", payload.UserID).Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			http.Error(writer, "User not found", http.StatusNotFound)
			return
		}

		/*	transactions, err := getTransactions(db, payload.AccountID, payload.Quarter)
			if err != nil {
				http.Error(writer, "Error fetching transactions", http.StatusInternalServerError)
				return
			}
			totalBalance, numTransactions, avgDebit, avgCredit := calculateSummary(transactions)
			summaryMessage := fmt.Sprintf(
				"Saldo total: %.2f\nNúmero de transacciones en julio: %d\nNúmero de transacciones en agosto: %d\nMonto promedio de débito: %.2f\nMonto promedio de crédito: %.2f",
				totalBalance, numTransactions["July"], numTransactions["August"], avgDebit, avgCredit)
	*/
	var totalBalance, numTransactions, avgDebit, avgCredit float64

	summaryMessage := fmt.Sprintf(
		"Saldo total: %.2f\nNúmero de transacciones en julio: %d\nNúmero de transacciones en agosto: %d\nMonto promedio de débito: %.2f\nMonto promedio de crédito: %.2f",
		totalBalance, numTransactions, numTransactions, avgDebit, avgCredit)

	// Respond with the summary message
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(map[string]string{
		"message": summaryMessage,
	})
	log.Printf("Starting lambda summary - ", "method: ", request.Method)
}

func getTransactions(db *sql.DB, accountID int, quarter string) ([]Transaction, error) {
	rows, err := db.Query("SELECT id, date, transaction FROM transactions WHERE account_id = $1", accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []Transaction
	for rows.Next() {
		var t Transaction
		err := rows.Scan(&t.ID, &t.Date, &t.Transaction)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}

	return transactions, nil
}

func calculateSummary(transactions []Transaction) (float64, map[string]int, float64, float64) {
	totalBalance := 0.0
	numTransactions := map[string]int{"July": 0, "August": 0}
	totalDebit := 0.0
	totalCredit := 0.0
	numDebit := 0
	numCredit := 0

	for _, t := range transactions {
		totalBalance += t.Transaction
		month := getMonthFromDateString(t.Date)
		numTransactions[month]++
		if t.Transaction < 0 {
			totalDebit += t.Transaction
			numDebit++
		} else {
			totalCredit += t.Transaction
			numCredit++
		}
	}

	avgDebit := totalDebit / float64(numDebit)
	avgCredit := totalCredit / float64(numCredit)
	return totalBalance, numTransactions, avgDebit, avgCredit
}

func getMonthFromDateString(date string) string {
	// Assuming date format is mm/dd
	parts := strings.Split(date, "/")
	month, _ := strconv.Atoi(parts[0])
	if month == 7 {
		return "July"
	} else if month == 8 {
		return "August"
	}
	return ""
}

func sendEmail(to, subject, body string) {
	from := "your-email@example.com"
	password := "your-email-password"
	smtpHost := "smtp.example.com"
	smtpPort := "587"

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	auth := smtp.PlainAuth("", from, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(msg))
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	//	lambda.Start(LambdaHandler)
	db, err := database.NewPostgresConnection()
	if err != nil {
		fmt.Printf("Failed to connect to the database: %v\n", err)
		return
	}
	defer db.Close()

	summaryRepository := repository.NewSummaryRepository(db)
	summaryHandler := handler.NewSummaryHandler(&summaryRepository)

	http.HandleFunc("/", health)
	http.HandleFunc("/summary", summaryHandler.Summary)
	http.ListenAndServe(":3000", nil)
}
