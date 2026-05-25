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
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockS3Client struct {
	mock.Mock
}

func (m *mockS3Client) GetObject(ctx context.Context, params *awss3.GetObjectInput, optFns ...func(*awss3.Options)) (*awss3.GetObjectOutput, error) {
	args := m.Called(ctx, params)

	if out, ok := args.Get(0).(*awss3.GetObjectOutput); ok {
		return out, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockS3Client) PutObject(ctx context.Context, params *awss3.PutObjectInput, optFns ...func(*awss3.Options)) (*awss3.PutObjectOutput, error) {
	args := m.Called(ctx, params)
	if out, ok := args.Get(0).(*awss3.PutObjectOutput); ok {
		return out, args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *mockS3Client) DeleteObject(ctx context.Context, params *awss3.DeleteObjectInput, optFns ...func(*awss3.Options)) (*awss3.DeleteObjectOutput, error) {
	args := m.Called(ctx, params)
	if out, ok := args.Get(0).(*awss3.DeleteObjectOutput); ok {
		return out, args.Error(1)
	}

	return nil, args.Error(1)
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

const (
	testBucket = "test-bucket"
	testRegion = "us-east-1"
	testKey    = "customers/123/avatar/abc.jpg"
)

func TestUpload(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		key         string
		contentType string
		data        []byte
		mockFn      func(c *mockS3Client)
		wantErr     bool
	}{
		{
			name:        "success",
			key:         testKey,
			contentType: "image/jpeg",
			data:        []byte("image data"),
			mockFn: func(c *mockS3Client) {
				c.On("PutObject", mock.Anything, mock.Anything).
					Return(&awss3.PutObjectOutput{}, nil).Once()
			},
			wantErr: false,
		},
		{
			name:        "error - empty key",
			key:         "",
			contentType: "image/jpeg",
			data:        []byte("image data"),
			mockFn:      func(c *mockS3Client) {},
			wantErr:     true,
		},
		{
			name:        "error - invalid key characters",
			key:         "customers/123/avatar/file name.jpg",
			contentType: "image/jpeg",
			data:        []byte("image data"),
			mockFn:      func(c *mockS3Client) {},
			wantErr:     true,
		},
		{
			name:        "error - s3 client fails",
			key:         testKey,
			contentType: "image/jpeg",
			data:        []byte("image data"),
			mockFn: func(c *mockS3Client) {
				c.On("PutObject", mock.Anything, mock.Anything).
					Return(&awss3.PutObjectOutput{}, errors.New("s3 error")).Once()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			client := new(mockS3Client)
			tt.mockFn(client)

			fakePresign := newFakePresignClient()

			storage := NewS3Storage(client, testBucket, "", testRegion, fakePresign, 15*time.Minute)
			err := storage.Upload(context.Background(), tt.key, tt.contentType, bytes.NewReader(tt.data))

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			client.AssertExpectations(t)
		})
	}
}

func TestGetPresignedURL(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	mockClient := new(mockS3Client)
	fakePresign := newFakePresignClient()

	storage := NewS3Storage(mockClient, testBucket, "", testRegion, fakePresign, 15*time.Minute)

	t.Run("Success generation", func(t *testing.T) {
		key := "posts/123/img.jpg"

		url, err := storage.GetPresignedURL(ctx, key)

		require.NoError(t, err)
		require.NotEmpty(t, url)
		require.True(t, strings.Contains(url, testBucket), "URL should contain bucket name")
		require.True(t, strings.Contains(url, key), "URL should contain object key")
		require.True(t, strings.Contains(url, "X-Amz-Signature"), "URL should contain AWS signature parameters")
	})

	t.Run("Empty key validation", func(t *testing.T) {
		url, err := storage.GetPresignedURL(ctx, "")

		require.Error(t, err)
		require.Empty(t, url)
	})

	mockClient.AssertNotCalled(t, "PutObject")
	mockClient.AssertNotCalled(t, "DeleteObject")
}

func TestDelete(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		key     string
		mockFn  func(c *mockS3Client)
		wantErr bool
	}{
		{
			name: "success",
			key:  testKey,
			mockFn: func(c *mockS3Client) {
				c.On("DeleteObject", mock.Anything, mock.Anything).
					Return(&awss3.DeleteObjectOutput{}, nil).Once()
			},
			wantErr: false,
		},
		{
			name:    "error - empty key",
			key:     "",
			mockFn:  func(c *mockS3Client) {},
			wantErr: true,
		},
		{
			name:    "error - invalid key characters",
			key:     "customers/123/avatar/file name.jpg",
			mockFn:  func(c *mockS3Client) {},
			wantErr: true,
		},
		{
			name: "error - s3 client fails",
			key:  testKey,
			mockFn: func(c *mockS3Client) {
				c.On("DeleteObject", mock.Anything, mock.Anything).
					Return(&awss3.DeleteObjectOutput{}, errors.New("s3 error")).Once()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			client := new(mockS3Client)
			tt.mockFn(client)

			fakePresign := newFakePresignClient()

			storage := NewS3Storage(client, testBucket, "", testRegion, fakePresign, 15*time.Minute)
			err := storage.Delete(context.Background(), tt.key)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			client.AssertExpectations(t)
		})
	}
}

func TestBuildURL(t *testing.T) {
	t.Parallel()

	t.Run("aws s3 url", func(t *testing.T) {
		t.Parallel()
		storage := NewS3Storage(nil, testBucket, "", testRegion, nil, 0)
		url := storage.BuildURL(testKey)
		require.Equal(t, "https://test-bucket.s3.us-east-1.amazonaws.com/customers/123/avatar/abc.jpg", url)
	})

	t.Run("garage custom endpoint url", func(t *testing.T) {
		t.Parallel()
		storage := NewS3Storage(nil, testBucket, "http://localhost:4300", "", nil, 0)
		url := storage.BuildURL(testKey)
		require.Equal(t, "http://localhost:4300/test-bucket/customers/123/avatar/abc.jpg", url)
	})
}

func TestGet(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		key     string
		mockFn  func(c *mockS3Client)
		wantErr bool
	}{
		{
			name: "success",
			key:  testKey,
			mockFn: func(c *mockS3Client) {
				c.On("GetObject", mock.Anything, mock.Anything).
					Return(
						&awss3.GetObjectOutput{
							Body: io.NopCloser(
								strings.NewReader("image data"),
							),
						},
						nil,
					).Once()
			},
			wantErr: false,
		},
		{
			name:    "error - empty key",
			key:     "",
			mockFn:  func(c *mockS3Client) {},
			wantErr: true,
		},
		{
			name:    "error - invalid key",
			key:     "invalid key.jpg",
			mockFn:  func(c *mockS3Client) {},
			wantErr: true,
		},
		{
			name: "error - s3 failure",
			key:  testKey,
			mockFn: func(c *mockS3Client) {
				c.On("GetObject", mock.Anything, mock.Anything).
					Return(
						(*awss3.GetObjectOutput)(nil),
						errors.New("s3 error"),
					).Once()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client := new(mockS3Client)

			tt.mockFn(client)

			fakePresign := newFakePresignClient()

			storage := NewS3Storage(
				client,
				testBucket,
				"",
				testRegion,
				fakePresign,
				15*time.Minute,
			)

			reader, err := storage.Get(
				context.Background(),
				tt.key,
			)

			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, reader)
			} else {
				require.NoError(t, err)
				require.NotNil(t, reader)

				_ = reader.Close()
			}

			client.AssertExpectations(t)
		})
	}
}
