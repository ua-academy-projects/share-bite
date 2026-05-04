package env

import (
	"os"
	"time"
)

type S3StorageConfig struct {
	StorageEndpoint     string
	StorageRegion       string
	StorageAccessKey    string
	StorageSecretKey    string
	StorageBucket       string
	StorageUsePathStyle bool
	StoragePresignTTL   time.Duration
}

func NewS3StorageConfig() (*S3StorageConfig, error) {
	ttl := convertTTLIntoDuration(os.Getenv("S3_PRESIGN_URL_TTL"))

	return &S3StorageConfig{
		StorageEndpoint:     os.Getenv("S3_ENDPOINT"),
		StorageRegion:       os.Getenv("S3_REGION"),
		StorageAccessKey:    os.Getenv("S3_ACCESS_KEY"),
		StorageSecretKey:    os.Getenv("S3_SECRET_KEY"),
		StorageBucket:       os.Getenv("S3_BUCKET"),
		StorageUsePathStyle: os.Getenv("S3_USE_PATH_STYLE") == "true",
		StoragePresignTTL:   ttl,
	}, nil
}

func (c *S3StorageConfig) Endpoint() string {
	return c.StorageEndpoint
}

func (c *S3StorageConfig) Region() string {
	return c.StorageRegion
}

func (c *S3StorageConfig) AccessKey() string {
	return c.StorageAccessKey
}

func (c *S3StorageConfig) SecretKey() string {
	return c.StorageSecretKey
}

func (c *S3StorageConfig) Bucket() string {
	return c.StorageBucket
}

func (c *S3StorageConfig) UsePathStyle() bool {
	return c.StorageUsePathStyle
}

func (c *S3StorageConfig) PresignTTL() time.Duration {
	return c.StoragePresignTTL
}

func convertTTLIntoDuration(str string) time.Duration {
	tempDur := 15 * time.Minute
	if str == "" {
		return tempDur
	}
	duration, err := time.ParseDuration(str)
	if err != nil {
		return tempDur
	}
	return duration
}
