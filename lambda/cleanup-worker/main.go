package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ua-academy-projects/share-bite/internal/cleanup/repository"
	"github.com/ua-academy-projects/share-bite/internal/cleanup/service"
	"github.com/ua-academy-projects/share-bite/internal/config"
	"github.com/ua-academy-projects/share-bite/pkg/database/pg"
	"github.com/ua-academy-projects/share-bite/pkg/logger"
)

type CleanupRequest struct {
	JobNames []string `json:"job_names"`
}

type CleanupResponse struct {
	Success      bool                     `json:"success"`
	TotalFound   int64                    `json:"total_found"`
	TotalDeleted int64                    `json:"total_deleted"`
	DryRun       bool                     `json:"dry_run"`
	Duration     string                   `json:"duration"`
	Results      []map[string]interface{} `json:"results"`
	Error        string                   `json:"error,omitempty"`
}

func LambdaHandler(ctx context.Context, request CleanupRequest) (CleanupResponse, error) {
	if err := config.Load(); err != nil {
		log.Printf("Warning: could not load config: %v", err)
	}

	cfg := config.Config()

	logger.Info(ctx, "Lambda: Connecting to database")
	client, err := pg.NewClient(ctx, cfg.Postgres.Dsn())
	if err != nil {
		errMsg := fmt.Sprintf("failed to create database client: %v", err)
		logger.Error(ctx, errMsg)
		return CleanupResponse{
			Success: false,
			Error:   errMsg,
			DryRun:  cfg.Cleanup.IsDryRun(),
		}, err
	}
	defer client.Close()

	if err := client.DB().Ping(ctx); err != nil {
		errMsg := fmt.Sprintf("failed to ping database: %v", err)
		logger.Error(ctx, errMsg)
		return CleanupResponse{
			Success: false,
			Error:   errMsg,
			DryRun:  cfg.Cleanup.IsDryRun(),
		}, err
	}

	logger.Info(ctx, "Lambda: Successfully connected to database")

	cleanupRepo := repository.NewPostgresCleanupRepository(client)
	cleanupSvc := service.NewCleanupService(cleanupRepo)

	logger.Info(ctx, fmt.Sprintf("Lambda: Starting cleanup with retention period: %v, batch size: %d, dry-run: %v",
		cfg.Cleanup.GetRetentionPeriod(),
		cfg.Cleanup.GetBatchSize(),
		cfg.Cleanup.IsDryRun(),
	))

	results, err := cleanupSvc.RunAllCleanups(
		ctx,
		cfg.Cleanup.GetRetentionPeriod(),
		cfg.Cleanup.GetBatchSize(),
		cfg.Cleanup.IsDryRun(),
	)

	if err != nil {
		errMsg := fmt.Sprintf("cleanup failed: %v", err)
		logger.Error(ctx, errMsg)
		return CleanupResponse{
			Success: false,
			Error:   errMsg,
			DryRun:  cfg.Cleanup.IsDryRun(),
		}, err
	}

	response := CleanupResponse{
		Success: true,
		DryRun:  cfg.Cleanup.IsDryRun(),
	}

	totalDuration := 0
	for _, result := range results {
		response.TotalFound += result.RecordsFound
		response.TotalDeleted += result.RecordsDeleted
		totalDuration += int(result.Duration().Seconds())

		resultMap := map[string]interface{}{
			"name":            result.Name,
			"records_found":   result.RecordsFound,
			"records_deleted": result.RecordsDeleted,
			"dry_run":         result.DryRun,
			"duration_ms":     result.Duration().Milliseconds(),
		}

		if len(result.Errors) > 0 {
			resultMap["errors"] = result.Errors
		}

		response.Results = append(response.Results, resultMap)
	}

	response.Duration = fmt.Sprintf("%ds", totalDuration)

	logger.Info(ctx, fmt.Sprintf("Lambda: Cleanup completed. Found: %d, Deleted: %d, Duration: %s",
		response.TotalFound, response.TotalDeleted, response.Duration))

	return response, nil
}

func main() {
	lambda.Start(LambdaHandler)
}
