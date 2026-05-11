package storage

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	awscfg "github.com/aws/aws-sdk-go-v2/config"
	awscred "github.com/aws/aws-sdk-go-v2/credentials"
	s3sdk "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/ua-academy-projects/share-bite/internal/config"
	"github.com/ua-academy-projects/share-bite/internal/storage/s3"
)

func NewStorageClient(ctx context.Context, cfg config.Storage) (ObjectStorage, error) {
	loaderOpts := []func(*awscfg.LoadOptions) error{
		awscfg.WithRegion(cfg.Region()),
	}

	if cfg.AccessKey() != "" && cfg.SecretKey() != "" {
		loaderOpts = append(loaderOpts, awscfg.WithCredentialsProvider(
			awscred.NewStaticCredentialsProvider(cfg.AccessKey(), cfg.SecretKey(), ""),
		))
	}

	awsCfg, err := awscfg.LoadDefaultConfig(ctx, loaderOpts...)
	if err != nil {
		return nil, err
	}

	s3Client := s3sdk.NewFromConfig(awsCfg, func(o *s3sdk.Options) {
		o.UsePathStyle = cfg.UsePathStyle()

		// BaseEndpoint is set when cfg.Endpoint is passed
		// This option makes possible using different solutions (e.g. Garage, MinIO) instead of original AWS S3
		if endpoint := cfg.Endpoint(); len(endpoint) > 0 {
			o.BaseEndpoint = aws.String(endpoint)
		}
	})

	presignClient := s3sdk.NewPresignClient(s3Client)

	ttl := cfg.PresignTTL()
	region := cfg.Region()

	return s3.NewS3Storage(s3Client, cfg.Bucket(), cfg.Endpoint(), cfg.Region(), presignClient, cfg.PresignTTL()), nil
}
