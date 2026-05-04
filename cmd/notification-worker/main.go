package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ua-academy-projects/share-bite/internal/notification/worker"
)

func main() {
	validator := worker.NewDefaultValidator()
	processor := &worker.NoOpProcessor{}
	h := worker.New(validator, processor)

	lambda.Start(h.HandleBatch)
}
