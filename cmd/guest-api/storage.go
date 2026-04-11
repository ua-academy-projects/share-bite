package main

import (
	"context"
	"github.com/ua-academy-projects/share-bite/internal/config"

	aws "github.com/aws/aws-sdk-go-v2/aws"
	awscfg "github.com/aws/aws-sdk-go-v2/config"
	awscred "github.com/aws/aws-sdk-go-v2/credentials"
	s3sdk "github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/ua-academy-projects/share-bite/internal/storage/s3"
)

func newStorageClient(ctx context.Context, cfg config.Storage) (*s3.S3Storage, error) {
	if cfg == nil || cfg.Bucket() == "" {
		return nil, nil
	}

	loaderOpts := []func(*awscfg.LoadOptions) error{
		awscfg.WithRegion(cfg.Region()),
	}

	if cfg.AccessKey() != "" && cfg.SecretKey() != "" {
		loaderOpts = append(loaderOpts, awscfg.WithCredentialsProvider(
			awscred.NewStaticCredentialsProvider(cfg.AccessKey(), cfg.SecretKey(), ""),
		))
	}

	if cfg.Endpoint() != "" {
		loaderOpts = append(loaderOpts, awscfg.WithEndpointResolverWithOptions(
			aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL:               cfg.Endpoint(),
					HostnameImmutable: true,
				}, nil
			}),
		))
	}

	awsCfg, err := awscfg.LoadDefaultConfig(ctx, loaderOpts...)
	if err != nil {
		return nil, err
	}

	s3Client := s3sdk.NewFromConfig(awsCfg, func(o *s3sdk.Options) {
		o.UsePathStyle = cfg.UsePathStyle()
	})

	return s3.NewS3Storage(s3Client, cfg.Bucket(), cfg.Endpoint()), nil
}
