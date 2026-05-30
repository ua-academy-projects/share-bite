package env

import (
	"time"
)

type s3StorageConfig struct {
	StorageEndpoint     string        `env:"S3_ENDPOINT"`
	StorageRegion       string        `env:"S3_REGION"`
	StorageAccessKey    string        `env:"S3_ACCESS_KEY"`
	StorageSecretKey    string        `env:"S3_SECRET_KEY"`
	StorageBucket       string        `env:"S3_BUCKET"`
	StorageUsePathStyle bool          `env:"S3_USE_PATH_STYLE" envDefault:"false"`
	StoragePresignTTL   time.Duration `env:"S3_PRESIGN_URL_TTL" envDefault:"15m"`
}

func NewS3StorageConfig(opts ...Options) (*s3StorageConfig, error) {
	config := new(s3StorageConfig)
	if err := Parse(config, opts...); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *s3StorageConfig) Endpoint() string {
	return c.StorageEndpoint
}

func (c *s3StorageConfig) Region() string {
	return c.StorageRegion
}

func (c *s3StorageConfig) AccessKey() string {
	return c.StorageAccessKey
}

func (c *s3StorageConfig) SecretKey() string {
	return c.StorageSecretKey
}

func (c *s3StorageConfig) Bucket() string {
	return c.StorageBucket
}

func (c *s3StorageConfig) UsePathStyle() bool {
	return c.StorageUsePathStyle
}

func (c *s3StorageConfig) PresignTTL() time.Duration {
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
