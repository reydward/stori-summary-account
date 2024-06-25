package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"load-data/internal/model"
	"load-data/internal/repository"
	"load-data/internal/service"
	"net/http"
	"strconv"
	"strings"
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
	// Getting the request params
	var payload model.RequestPayload

	accountIDStr := request.FormValue("accountId")
	if accountIDStr == "" {
		http.Error(writer, "Falta el ID de cuenta", http.StatusBadRequest)
		return
	}

	err := request.ParseMultipartForm(10 << 20) // 10 MB file size limit
	if err != nil {
		http.Error(writer, "Could not parse multipart form", http.StatusBadRequest)
		return
	}

	file, fileHeader, err := request.FormFile("file")
	if err != nil {
		http.Error(writer, "Error getting the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validating the file type
	if !strings.HasPrefix(fileHeader.Header.Get("Content-Type"), "text/csv") {
		http.Error(writer, "Invalid file type", http.StatusBadRequest)
		return
	}

	// Reading the file content
	var fileContent bytes.Buffer
	if _, err := io.Copy(&fileContent, file); err != nil {
		http.Error(writer, "Could not read file content", http.StatusInternalServerError)
		return
	}

	// Setting the payload
	accountID, err := strconv.Atoi(accountIDStr)
	if err != nil {
		http.Error(writer, "Invalid account ID", http.StatusBadRequest)
		return
	}

	payload.AccountID = accountID
	payload.FileName = fileHeader.Filename
	payload.File = fileContent

	fmt.Println("Payload: ", payload)

	// Processing the transaction
	responseMessage, err := service.ProcessTransactions(payload, h.repo)
	if err != nil {
		http.Error(writer, responseMessage, http.StatusInternalServerError)
		return
	}

	//Setting the response
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(struct{ Message string }{responseMessage})
}
