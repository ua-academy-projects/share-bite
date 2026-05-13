package main

import (
	"context"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ua-academy-projects/share-bite/internal/config"
	notificationworker "github.com/ua-academy-projects/share-bite/internal/notification/worker"
	"github.com/ua-academy-projects/share-bite/pkg/email"
	"github.com/ua-academy-projects/share-bite/pkg/logger"
	"github.com/ua-academy-projects/share-bite/pkg/notification"
)

func main() {

	ctx := context.Background()

	secretName := os.Getenv("AWS_LAMBDA_SECRETS_NAME")

	if secretName != "" {
		if err := config.LoadFromAWSSecrets(ctx, secretName); err != nil {
			logger.Fatal(ctx, "failed to load secrets from AWS", err)
		}
	}

	if err := config.Load(); err != nil {
		if strings.HasPrefix(err.Error(), "missing required environment variables:") {
			logger.FatalKV(ctx, "config load failed", "missing_envs", err.Error())
		} else {
			logger.Fatal(ctx, "config load:", err)
		}
	}

	var emailSender email.Sender
	if config.Config().App.IsProd() {
		emailSender = email.NewResendSender(
			config.Config().Email.ResendAPIKeyValue(),
			config.Config().Email.ResendFromEmailValue(),
		)
	} else {
		emailSender = email.NewFakeSender()
	}

	validator := notificationworker.NewDefaultValidator(notification.RegistrationConfirmed)

	emailProcessor := notificationworker.NewEmailProcessor(emailSender)

	h := notificationworker.New(validator, emailProcessor)

	lambda.Start(h.HandleBatch)
}
