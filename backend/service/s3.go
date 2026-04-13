package service

import (
	"context"
	"log"
	"mime/multipart"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var S3Client *s3.Client

func defaultBucketName() string {
	return os.Getenv("AWS_BUCKET_NAME")
}

func InitS3() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal("unable to load AWS config:", err)
	}

	S3Client = s3.NewFromConfig(cfg)
}

func GenerateSignedURL(bucket string, key string) (string, error) {

	presignClient := s3.NewPresignClient(S3Client)

	req, err := presignClient.PresignGetObject(context.TODO(),
		&s3.GetObjectInput{
			Bucket: &bucket,
			Key:    &key,
		},
		s3.WithPresignExpires(5*time.Minute),
	)

	if err != nil {
		return "", err
	}

	return req.URL, nil
}

func UploadToS3(file *multipart.FileHeader, key string, bucket string) error {

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	contentType := file.Header.Get("Content-Type")

	input := &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &key,
		Body:   src,
	}

	if contentType != "" {
		input.ContentType = aws.String(contentType)
	}

	_, err = S3Client.PutObject(context.TODO(), input)

	return err
}

func DeleteFileFromS3(storageKey string) error {

	_, err := S3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(defaultBucketName()),
		Key:    aws.String(storageKey),
	})

	return err
}
