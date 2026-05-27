package s3

import (
	"context"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"

	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	// safeKeyRegexp matches S3 object keys using only recommended safe characters.
	// This is intentionally strict — keys are only programmatically-generated S3 key characters.
	safeKeyRegexp = regexp.MustCompile(`^[a-zA-Z0-9/\-_.]+$`)
)

type s3Client interface {
	PutObject(ctx context.Context, params *awss3.PutObjectInput, optFns ...func(*awss3.Options)) (*awss3.PutObjectOutput, error)
	DeleteObject(ctx context.Context, params *awss3.DeleteObjectInput, optFns ...func(*awss3.Options)) (*awss3.DeleteObjectOutput, error)
	GetObject(ctx context.Context, params *awss3.GetObjectInput, optFns ...func(*awss3.Options)) (*awss3.GetObjectOutput, error)
}

type S3Storage struct {
	client     s3Client
	bucketName string
	endpoint   string
	region     string

	presignClient *awss3.PresignClient
	urlTTL        time.Duration
}

func NewS3Storage(client s3Client, bucketName string, endpoint string, region string, presignClient *awss3.PresignClient, ttl time.Duration) *S3Storage {
	return &S3Storage{
		client:     client,
		bucketName: bucketName,
		endpoint:   endpoint,
		region:     region,

		presignClient: presignClient,
		urlTTL:        ttl,
	}
}

func (s *S3Storage) Upload(ctx context.Context, key string, contentType string, data io.Reader) error {
	if err := validateKey(key); err != nil {
		return err
	}

	_, err := s.client.PutObject(ctx, &awss3.PutObjectInput{
		Bucket:      &s.bucketName,
		Key:         &key,
		Body:        data,
		ContentType: &contentType,
	})
	if err != nil {
		return fmt.Errorf("failed to upload object to storage: %w", err)
	}

	return nil
}

func (s *S3Storage) Delete(ctx context.Context, key string) error {
	if err := validateKey(key); err != nil {
		return err
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

// BuildURL constructs a public URL for the given object key.
// When a custom endpoint is configured (e.g. Garage, MinIO for local development),
// it uses path-style URL format. Otherwise it falls back to the standard AWS S3 URL format.
func (s *S3Storage) BuildURL(key string) string {
	if len(s.endpoint) > 0 {
		return fmt.Sprintf("%s/%s/%s", s.endpoint, s.bucketName, key)
	}

	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.bucketName, s.region, key)
}

func (s *S3Storage) GetPresignedURL(ctx context.Context, key string) (string, error) {
	if key == "" {
		return "", fmt.Errorf("object key is required")
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

func (s *S3Storage) Get(ctx context.Context, key string) (io.ReadCloser, error) {
	if err := validateKey(key); err != nil {
		return nil, err
	}

	result, err := s.client.GetObject(
		ctx,
		&awss3.GetObjectInput{
			Bucket: &s.bucketName,
			Key:    &key,
		},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get object from storage: %w",
			err,
		)
	}

	return result.Body, nil
}

// validateKey simply checks whether key is ok to be used
func validateKey(key string) error {
	if key == "" {
		return fmt.Errorf("object key is required")
	}
	if !safeKeyRegexp.MatchString(key) {
		return fmt.Errorf("invalid object key is provided")
	}

	for s := range strings.SplitSeq(key, "/") {
		if s == "" || s == "." || s == ".." {
			return fmt.Errorf("invalid object key is provided")
		}
	}

	return nil
}
