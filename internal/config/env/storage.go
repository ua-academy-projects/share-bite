package env

import "os"

type S3StorageConfig struct {
	StorageEndpoint     string
	StorageRegion       string
	StorageAccessKey    string
	StorageSecretKey    string
	StorageBucket       string
	StorageUsePathStyle bool
}

func NewS3StorageConfig() (*S3StorageConfig, error) {
	return &S3StorageConfig{
		StorageEndpoint:     os.Getenv("S3_ENDPOINT"),
		StorageRegion:       os.Getenv("S3_REGION"),
		StorageAccessKey:    os.Getenv("S3_ACCESS_KEY"),
		StorageSecretKey:    os.Getenv("S3_SECRET_KEY"),
		StorageBucket:       os.Getenv("S3_BUCKET"),
		StorageUsePathStyle: os.Getenv("S3_USE_PATH_STYLE") == "true",
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
