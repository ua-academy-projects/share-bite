package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ua-academy-projects/share-bite/internal/cleanup/repository"
	"github.com/ua-academy-projects/share-bite/internal/cleanup/service"
	"github.com/ua-academy-projects/share-bite/internal/config"
	"github.com/ua-academy-projects/share-bite/pkg/database/pg"
	"github.com/ua-academy-projects/share-bite/pkg/logger"
)

func main() {
	dryRun := flag.Bool("dry-run", false, "Enable dry-run mode (log what would be deleted without deleting)")
	help := flag.Bool("help", false, "Show help message")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `
ShareBite Cleanup Worker

Usage: cleanup-worker [flags]

Flags:
  -dry-run    Enable dry-run mode (logs what would be deleted without deleting)
  -help       Show this help message

Environment Variables:
  CLEANUP_RETENTION_DAYS      Number of days to retain records (default: 30)
  CLEANUP_BATCH_SIZE          Number of records to delete per batch (default: 100)
  CLEANUP_DRY_RUN             Set to "true" to enable dry-run mode (can be overridden by -dry-run flag)
  CLEANUP_SCHEDULE_ENABLED    Set to "true" to enable scheduling (used by Lambda handler)

Examples:
  # Run cleanup with defaults
  cleanup-worker

  # Run in dry-run mode
  cleanup-worker -dry-run

  # Run with environment variables
  CLEANUP_RETENTION_DAYS=60 CLEANUP_BATCH_SIZE=500 cleanup-worker

`)
	}

	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	ctx := context.Background()

	if err := config.Load(".env"); err != nil {
		log.Printf("Warning: could not load .env file: %v", err)
	}

	cfg := config.Config()

	if cfg.App.IsProd() {
		logger.Info(ctx, "Running in production mode")
	} else {
		logger.Info(ctx, "Running in development mode")
	}

	logger.Info(ctx, "Connecting to database...")

	client, err := pg.NewClient(ctx, cfg.Postgres.Dsn())
	if err != nil {
		logger.Fatal(ctx, fmt.Sprintf("failed to create database client: %v", err))
	}
	defer client.Close()

	if err := client.DB().Ping(ctx); err != nil {
		logger.Fatal(ctx, fmt.Sprintf("failed to ping database: %v", err))
	}
	logger.Info(ctx, "Successfully connected to database")

	cleanupRepo := repository.NewPostgresCleanupRepository(client)
	cleanupSvc := service.NewCleanupService(cleanupRepo)

	dryRunMode := *dryRun || cfg.Cleanup.IsDryRun()

	logger.Info(ctx, fmt.Sprintf("Starting cleanup with retention period: %v, batch size: %d, dry-run: %v",
		cfg.Cleanup.GetRetentionPeriod(),
		cfg.Cleanup.GetBatchSize(),
		dryRunMode,
	))

	results, err := cleanupSvc.RunAllCleanups(
		ctx,
		cfg.Cleanup.GetRetentionPeriod(),
		cfg.Cleanup.GetBatchSize(),
		dryRunMode,
	)

	if err != nil {
		logger.Error(ctx, fmt.Sprintf("cleanup failed: %v", err))
		os.Exit(1)
	}

	fmt.Println("\n=== Cleanup Results ===")
	totalFound := int64(0)
	totalDeleted := int64(0)

	for _, result := range results {
		fmt.Printf("\nJob: %s\n", result.Name)
		fmt.Printf("  Records Found:  %d\n", result.RecordsFound)
		if result.DryRun {
			fmt.Printf("  Would Delete:   %d (DRY-RUN)\n", result.RecordsDeleted)
		} else {
			fmt.Printf("  Records Deleted: %d\n", result.RecordsDeleted)
		}
		fmt.Printf("  Duration:       %v\n", result.Duration())
		if len(result.Errors) > 0 {
			fmt.Printf("  Errors:\n")
			for _, errMsg := range result.Errors {
				fmt.Printf("    - %s\n", errMsg)
			}
		}
		totalFound += result.RecordsFound
		totalDeleted += result.RecordsDeleted
	}

	fmt.Printf("\n=== Summary ===\n")
	if dryRunMode {
		fmt.Printf("Total Records Found:     %d\n", totalFound)
		fmt.Printf("Would Delete (DRY-RUN):  %d\n", totalDeleted)
		fmt.Printf("Mode: DRY-RUN (no changes made)\n")
	} else {
		fmt.Printf("Total Records Found:  %d\n", totalFound)
		fmt.Printf("Total Records Deleted: %d\n", totalDeleted)
	}
	fmt.Println()

	os.Exit(0)
}
