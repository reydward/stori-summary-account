package file

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io"
	"log"
	"strings"
)

type AWSService struct {
	S3Client *s3.Client
}

func (awsService AWSService) UploadFile(bucketName string, bucketKey string, fileContent bytes.Buffer) error {
	_, err := awsService.S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(bucketKey),
		Body:   strings.NewReader(fileContent.String()),
	})
	if err != nil {
		fmt.Printf("Error uploading file", err)
	}

	return err
}

func (awsService AWSService) GetFile(bucketName string, bucketKey string) io.ReadCloser {
	log.Printf("Getting the file from S3: %s with key: %s", bucketName, bucketKey)
	response, err := awsService.S3Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(bucketKey),
	})
	if err != nil {
		fmt.Errorf("Error reading the file from S3", err)
	}

	return response.Body
}

func (awsService AWSService) DeleteFile(bucketName string, bucketKey string) error {
	_, err := awsService.S3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(bucketKey),
	})
	if err != nil {
		return fmt.Errorf("Error deleting the file from S3 %v: %v", bucketKey, err)
	}

	return nil
}
