package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ua-academy-projects/share-bite/internal/config"
	notificationworker "github.com/ua-academy-projects/share-bite/internal/notification/worker"
	"github.com/ua-academy-projects/share-bite/pkg/email"
	"github.com/ua-academy-projects/share-bite/pkg/logger"
	"github.com/ua-academy-projects/share-bite/pkg/notification"
)

func main() {
	ctx := context.Background()

	if err := config.Load(".env"); err != nil {
		logger.Fatal(ctx, "config load:", err)
	}

	var emailSender email.Sender
	if config.Config().App.IsProd() {
		emailSender = email.NewResendSender(config.Config().Email.ResendAPIKeyValue())
	} else {
		emailSender = email.NewFakeSender()
	}

	validator := notificationworker.NewDefaultValidator(notification.RegistrationConfirmed)

	emailProcessor := notificationworker.NewEmailProcessor(emailSender)

	h := notificationworker.New(validator, emailProcessor)

	lambda.Start(h.HandleBatch)
}
