package cloudflare

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go/aws"
)

func InitR2Client() (*s3.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion("auto"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("YOUR_ACCESS_KEY", "YOUR_SECRET_KEY", "")))
	if err != nil {
		return nil, fmt.Errorf("unable to load AWS config: %v", err)
	}

	client := s3.NewFromConfig(cfg)
	return client, nil
}

func UploadFile(bucketName, fileName string, fileData io.Reader) (string, error) {
	client, err := InitR2Client()
	if err != nil {
		return "", err
	}

	// Prepare the S3 PutObject request
	_, err = client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
		Body:   fileData,
		ACL:    s3.ObjectCannedACLPublicRead, // Set the file to be publicly accessible
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %v", err)
	}

	// Construct the URL of the uploaded file
	url := fmt.Sprintf("https://%s.r2.cloudflarestorage.com/%s/%s", "YOUR_ACCOUNT_ID", bucketName, fileName)
	return url, nil
}
