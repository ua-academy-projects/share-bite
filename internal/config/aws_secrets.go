package config

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	awscfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/ua-academy-projects/share-bite/internal/config/env"
	"github.com/ua-academy-projects/share-bite/pkg/logger"
)

const notificationPrefixLambda = "NOTIFICATION_"

func LoadFromAWSSecrets(ctx context.Context, secretName string) error {
	awsConfig, err := awscfg.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := secretsmanager.NewFromConfig(awsConfig)

	logger.Infof(ctx, "Loading configuration from AWS Secrets Manager: %s", secretName)
	result, err := client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId: &secretName,
	})
	if err != nil {
		return fmt.Errorf("failed to get secret from AWS Secrets Manager: %w", err)
	}

	var secretMap map[string]string
	if err := json.Unmarshal([]byte(*result.SecretString), &secretMap); err != nil {
		return fmt.Errorf("failed to parse secret JSON: %w", err)
	}

	for key, value := range secretMap {
		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("failed to set environment variable %s: %w", key, err)
		}
	}

	logger.Infof(ctx, "Successfully loaded %d environment variables from AWS Secrets Manager", len(secretMap))

	appConfig, err := env.NewAppConfig()
	if err != nil {
		return fmt.Errorf("app config: %w", err)
	}

	postgresConfig, err := env.NewPostgresConfig()
	if err != nil {
		return fmt.Errorf("postgres config: %w", err)
	}

	redisConfig, err := env.NewRedisConfig()
	if err != nil {
		return fmt.Errorf("redis config: %w", err)
	}

	jwtTokenConfig, err := env.NewJwtTokenConfig()
	if err != nil {
		return fmt.Errorf("jwt token config: %w", err)
	}

	emailConfig, err := env.NewEmailConfig()
	if err != nil {
		return fmt.Errorf("email config: %w", err)
	}

	notificationHttpServerConfig, err := env.NewHttpServerConfig(notificationPrefixLambda)
	if err != nil {
		return fmt.Errorf("notification http server config: %w", err)
	}

	notificationSQSConfig, err := env.NewSQSConfig(notificationPrefixLambda)
	if err != nil {
		return fmt.Errorf("notification sqs config: %w", err)
	}

	cfg = &config{
		App:      appConfig,
		Postgres: postgresConfig,
		Redis:    redisConfig,
		JwtToken: jwtTokenConfig,
		Email:    emailConfig,

		NotificationHttpServer: notificationHttpServerConfig,
		NotificationSQS:        notificationSQSConfig,
	}

	return nil
}
