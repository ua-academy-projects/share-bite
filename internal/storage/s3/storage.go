package s3

import (
	"context"
	"fmt"
	"io"
	"time"

	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
)

type S3API interface {
	PutObject(ctx context.Context, params *awss3.PutObjectInput, optFns ...func(*awss3.Options)) (*awss3.PutObjectOutput, error)
	DeleteObject(ctx context.Context, params *awss3.DeleteObjectInput, optFns ...func(*awss3.Options)) (*awss3.DeleteObjectOutput, error)
}

type S3Storage struct {
	client        S3API
	bucketName    string
	endpoint      string
	presignClient *awss3.PresignClient
	urlTTL           time.Duration
}

func NewS3Storage(client S3API, bucketName string, endpoint string, presignClient *awss3.PresignClient, ttl time.Duration) *S3Storage {
	return &S3Storage{
		client:        client,
		bucketName:    bucketName,
		endpoint:      endpoint,
		presignClient: presignClient,
		urlTTL:        ttl,
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

// Deprecated: use GetPresignedURL instead and save object key to db either full image string.
// This method to be replaced with complete AWS integration
func (s *S3Storage) BuildURL(key string) string {
	return fmt.Sprintf("%s/%s/%s", s.endpoint, s.bucketName, key)
}

func (s *S3Storage) GetPresignedURL(ctx context.Context, key string) (string, error) {
	if key == "" {
		return "", apperror.BadRequest("object key is required")
	}

	req, err := s.presignClient.PresignGetObject(ctx, &awss3.GetObjectInput{
		Bucket: &s.bucketName,
		Key:    &key,
	}, func(po *awss3.PresignOptions) {
		po.Expires = s.urlTTL
	})

	if err != nil {
		return "", fmt.Errorf("failed to get object from S3 storage: %w", err)
	}

	return req.URL, nil
}
