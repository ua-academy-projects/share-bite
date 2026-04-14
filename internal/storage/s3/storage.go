package s3

import (
	"context"
	"fmt"
	"io"

	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
)

type S3Storage struct {
	client     *awss3.Client
	bucketName string
	endpoint   string
}

func NewS3Storage(client *awss3.Client, bucketName string, endpoint string) *S3Storage {
	return &S3Storage{
		client:     client,
		bucketName: bucketName,
		endpoint:   endpoint,
	}
}

func (s *S3Storage) Upload(
	ctx context.Context,
	key string,
	contentType string,
	data io.Reader,
) (string, error) {
	if key == "" {
		return "", apperror.BadRequest("object key is required")
	}

	if contentType == "" {
		return "", apperror.ErrUnsupportedImageType
	}

	if data == nil {
		return "", apperror.ErrImageRequired
	}

	_, err := s.client.PutObject(ctx, &awss3.PutObjectInput{
		Bucket:      &s.bucketName,
		Key:         &key,
		Body:        data,
		ContentType: &contentType,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload object to storage: %w", err)
	}

	return key, nil
}

func (s *S3Storage) Delete(ctx context.Context, key string) error {
	if key == "" {
		return apperror.BadRequest("object key is required")
	}

	_, err := s.client.DeleteObject(ctx, &awss3.DeleteObjectInput{
		Bucket: &s.bucketName,
		Key:    &key,
	})
	if err != nil {
		return fmt.Errorf("failed to delete object from storage: %w", err)
	}

	return nil
}

func (s *S3Storage) BuildURL(key string) string {
	return fmt.Sprintf("%s/%s/%s", s.endpoint, s.bucketName, key)
}