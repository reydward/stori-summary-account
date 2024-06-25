package handler

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/grokify/go-awslambda"
	"io"
	"load-data/internal/database"
	"load-data/internal/model"
	"load-data/internal/repository"
	"load-data/internal/service"
	"log"
	"mime/multipart"
	"strconv"
	"strings"
)

func LambdaHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Starting load data function")

	db, err := database.NewPostgresConnection()
	if err != nil {
		log.Printf("Failed to connect to the database: %v\n", err)
	}
	defer db.Close()

	var response string

	if request.HTTPMethod == "POST" {
		log.Printf("LambdaHandler.POST")

		//Getting the payload
		payload, err := getPayloadOld(request)
		if err != nil {
			log.Printf("Failed to get the payload: %v", err)
			return events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:       fmt.Sprintf(`{"error": "Invalid request payload: %v"}`, err),
			}, nil
		}

		log.Printf("Payload: %v", payload)
		repo := repository.NewLoadDataRepository(db)
		log.Printf("LambdaHandler.repo: %v", repo)

		// Processing the transaction
		response, err = service.ProcessTransactions(payload, repo)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       fmt.Sprintf(`{"error": "%s: %v"}`, response, err),
			}, nil
		}
	}

	var apigwresponse = &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       response,
	}

	apigwresponse.Headers = make(map[string]string)
	apigwresponse.Headers["Access-Control-Allow-Origin"] = "*"
	apigwresponse.Headers["Access-Control-Allow-Methods"] = "POST,OPTIONS"
	return *apigwresponse, nil
}

func getPayload(request events.APIGatewayProxyRequest) (model.RequestPayload, error) {
	fmt.Printf("getPayload.start")
	var payload model.RequestPayload
	var fileContent bytes.Buffer
	var fileName string
	var accountId string

	contentType := request.Headers["Content-Type"]
	if !strings.Contains(contentType, "multipart/form-data") {
		fmt.Printf("Content-Type must be multipart/form-data")
		return model.RequestPayload{}, errors.New("Content-Type must be multipart/form-data")
	}

	multipartReader, err := awslambda.NewReaderMultipart(request)
	if err != nil {
		return model.RequestPayload{}, err
	}

	part, err := multipartReader.NextPart()
	if err != nil {
		fmt.Printf("getPayload.multipartReader.NextPart().error: %v", err)
		return model.RequestPayload{}, err
	}
	for {
		if part.FormName() == "file" {
			fileName = part.FileName()
			_, err = io.Copy(&fileContent, part)
			if err != nil {
				fmt.Printf("getPayload.file.io.Copy.error: %v", err)
				return model.RequestPayload{}, err
			}
		}
		if part.FormName() == "accountId" {
			var buf bytes.Buffer
			_, err := io.Copy(&buf, part)
			if err != nil {
				fmt.Printf("getPayload.accountId.io.Copy.error: %v", err)
				return model.RequestPayload{}, err
			}
			accountId = buf.String()
		}

	}

	accountIDNumber, err := strconv.Atoi(accountId)
	if err != nil {
		fmt.Printf("getPayload.accountIDNumber.error: %v", err)
		return model.RequestPayload{}, err
	}

	payload.AccountID = accountIDNumber
	payload.FileName = fileName
	payload.File = fileContent
	fmt.Printf("getPayload.end")

	return payload, nil
}

// "Invalid request payload: multipart: NextPart: EOF"
func getPayloadOld(request events.APIGatewayProxyRequest) (model.RequestPayload, error) {
	fmt.Printf("getPayload.start")
	var payload model.RequestPayload
	var fileContent bytes.Buffer
	var accountId string

	body := bytes.NewReader([]byte(request.Body))
	reader := multipart.NewReader(body, request.Headers["Content-Type"])

	form, err := reader.ReadForm(32 << 20) // 32MB max memory
	if err != nil {
		return model.RequestPayload{}, err
	}
	defer form.RemoveAll()

	accountIds := form.Value["accountId"]
	if len(accountIds) == 0 {
		return model.RequestPayload{}, err
	}
	accountId = accountIds[0]

	files := form.File["file"]
	if len(files) == 0 {
		return model.RequestPayload{}, err
	}
	file := files[0]

	accountIDNumber, err := strconv.Atoi(accountId)
	if err != nil {
		fmt.Printf("getPayload.accountIDNumber.error: %v", err)
		return model.RequestPayload{}, err
	}

	payload.AccountID = accountIDNumber
	payload.FileName = file.Filename
	payload.File = fileContent
	fmt.Printf("getPayload.end")

	return payload, nil
}
