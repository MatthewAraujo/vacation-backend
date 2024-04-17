package main

import (
	"bytes"
	"context"
	"fmt"
	"log"

	configs "github.com/MatthewAraujo/vacation-backend/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Service struct {
	s3Client *s3.Client
	bucket   string
}

func NewR2Service() (*S3Service, error) {
	s3Config := configs.Envs.Cloudflare
	bucket := s3Config.BucketName
	account := s3Config.AccountID
	accessKey := s3Config.AccessKeyID
	secretKey := s3Config.AccessKeySecret

	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: fmt.Sprintf("https://%s.r2.cloudflarestorage.com", account),
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(r2Resolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		config.WithRegion("apac"),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Amazon S3 service client access key and so on
	s3Client := s3.NewFromConfig(cfg)

	return &S3Service{
		s3Client: s3Client,
		bucket:   bucket,
	}, nil

}
func (s *S3Service) UploadFileToR2(ctx context.Context, key string, file []byte) error {
	input := &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket), // Include the bucket name here
		Key:         aws.String(key),
		Body:        bytes.NewReader(file),
		ContentType: aws.String("image/jpeg"),
	}

	_, err := s.s3Client.PutObject(ctx, input)
	if err != nil {
		return err
	}

	return nil
}
