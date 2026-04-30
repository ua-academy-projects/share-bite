package s3

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"
)

type mockS3Client struct {
	putObjectErr    error
	deleteObjectErr error
}

func (m *mockS3Client) PutObject(ctx context.Context, params *awss3.PutObjectInput, optFns ...func(*awss3.Options)) (*awss3.PutObjectOutput, error) {
	if m.putObjectErr != nil {
		return nil, m.putObjectErr
	}
	return &awss3.PutObjectOutput{}, nil
}

func (m *mockS3Client) DeleteObject(ctx context.Context, params *awss3.DeleteObjectInput, optFns ...func(*awss3.Options)) (*awss3.DeleteObjectOutput, error) {
	if m.deleteObjectErr != nil {
		return nil, m.deleteObjectErr
	}
	return &awss3.DeleteObjectOutput{}, nil
}

func newFakePresignClient() *awss3.PresignClient {
	cfg := aws.Config{
		Region: "us-east-1",
		Credentials: aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
			return aws.Credentials{AccessKeyID: "fake", SecretAccessKey: "fake"}, nil
		}),
	}
	client := awss3.NewFromConfig(cfg)
	return awss3.NewPresignClient(client)
}

func TestS3Storage_Upload(t *testing.T) {
	ctx := context.Background()
	bucket := "test-bucket"
	fakePresign := newFakePresignClient()

	tests := []struct {
		name        string
		key         string
		contentType string
		data        []byte
		mockErr     error
		wantErr     bool
	}{
		{
			name:        "Success upload",
			key:         "posts/1/image.jpg",
			contentType: "image/jpeg",
			data:        []byte("fake-image-data"),
			mockErr:     nil,
			wantErr:     false,
		},
		{
			name:        "Empty key validation",
			key:         "",
			contentType: "image/jpeg",
			data:        []byte("fake-image-data"),
			mockErr:     nil,
			wantErr:     true, 
		},
		{
			name:        "Empty content type validation",
			key:         "posts/1/image.jpg",
			contentType: "",
			data:        []byte("fake-image-data"),
			mockErr:     nil,
			wantErr:     true, 
		},
		{
			name:        "Nil data validation",
			key:         "posts/1/image.jpg",
			contentType: "image/jpeg",
			data:        nil,
			mockErr:     nil,
			wantErr:     true,
		},
		{
			name:        "AWS SDK error",
			key:         "posts/1/image.jpg",
			contentType: "image/jpeg",
			data:        []byte("fake-image-data"),
			mockErr:     errors.New("aws internal error"),
			wantErr:     true, 
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockS3Client{putObjectErr: tt.mockErr}
			storage := NewS3Storage(mockClient, bucket, "http://localhost", fakePresign, 15*time.Minute)

			var reader io.Reader
			if tt.data != nil {
				reader = bytes.NewReader(tt.data)
			}

			gotKey, err := storage.Upload(ctx, tt.key, tt.contentType, reader)

			if (err != nil) != tt.wantErr {
				t.Errorf("Upload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && gotKey != tt.key {
				t.Errorf("Upload() gotKey = %v, want %v", gotKey, tt.key)
			}
		})
	}
}

func TestS3Storage_GetPresignedURL(t *testing.T) {
	ctx := context.Background()
	bucket := "test-bucket"
	
	mockClient := &mockS3Client{}
	fakePresign := newFakePresignClient()
	storage := NewS3Storage(mockClient, bucket, "http://localhost", fakePresign, 15*time.Minute)

	t.Run("Success generation", func(t *testing.T) {
		key := "posts/123/img.jpg"
		
		url, err := storage.GetPresignedURL(ctx, key)
		
		if err != nil {
			t.Fatalf("GetPresignedURL() unexpected error: %v", err)
		}
		
		if url == "" {
			t.Error("GetPresignedURL() returned empty URL")
		}

		if !strings.Contains(url, bucket) {
			t.Errorf("URL should contain bucket name, got: %s", url)
		}
		if !strings.Contains(url, key) {
			t.Errorf("URL should contain object key, got: %s", url)
		}
		if !strings.Contains(url, "X-Amz-Signature") {
			t.Errorf("URL should contain AWS signature parameters, got: %s", url)
		}
	})

	t.Run("Empty key validation", func(t *testing.T) {
		url, err := storage.GetPresignedURL(ctx, "")
		
		if err == nil {
			t.Error("GetPresignedURL() expected error for empty key, got nil")
		}
		if url != "" {
			t.Errorf("GetPresignedURL() expected empty URL for error, got %s", url)
		}
	})
}