package file_test

/*
import (
	"bytes"
	"io/ioutil"
	"load-data/internal/file"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUploadFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockS3 := mocks.NewMockS3API(ctrl)
	awsService := file.AWSService{S3Client: mockS3}

	bucketName := "test-bucket"
	bucketKey := "test-key"
	filePath := "test-file.txt"

	// Simular la apertura del archivo
	fileContent := []byte("this is a test file")
	ioutil.WriteFile(filePath, fileContent, 0644)
	defer os.Remove(filePath)

	file, _ := os.Open(filePath)
	defer file.Close()

	mockS3.EXPECT().PutObject(gomock.Any(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(bucketKey),
		Body:   file,
	}).Return(&s3.PutObjectOutput{}, nil)

	err := awsService.UploadFile(bucketName, bucketKey, filePath)
	assert.NoError(t, err)
}

func TestGetFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockS3 := mocks.NewMockS3API(ctrl)
	awsService := file.AWSService{S3Client: mockS3}

	bucketName := "test-bucket"
	bucketKey := "test-key"

	fileContent := []byte("this is a test file")
	readCloser := ioutil.NopCloser(bytes.NewReader(fileContent))

	mockS3.EXPECT().GetObject(gomock.Any(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(bucketKey),
	}).Return(&s3.GetObjectOutput{
		Body: readCloser,
	}, nil)

	body := awsService.GetFile(bucketName, bucketKey)
	defer body.Close()

	result, _ := ioutil.ReadAll(body)
	assert.Equal(t, fileContent, result)
}

func TestDeleteFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockS3 := mocks.NewMockS3API(ctrl)
	awsService := file.AWSService{S3Client: mockS3}

	bucketName := "test-bucket"
	bucketKey := "test-key"

	mockS3.EXPECT().DeleteObject(gomock.Any(), &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(bucketKey),
	}).Return(&s3.DeleteObjectOutput{}, nil)

	err := awsService.DeleteFile(bucketName, bucketKey)
	assert.NoError(t, err)
}
*/
