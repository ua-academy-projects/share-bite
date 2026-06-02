package env

import (
	"os"
	"strconv"
	"time"
)

type CleanupConfig struct {
	RetentionPeriod time.Duration
	BatchSize       int
	DryRun          bool
	ScheduleEnabled bool
}

func NewCleanupConfig() (*CleanupConfig, error) {
	retentionDaysStr := os.Getenv("CLEANUP_RETENTION_DAYS")
	retentionDays := 30
	if retentionDaysStr != "" {
		if days, err := strconv.Atoi(retentionDaysStr); err == nil {
			retentionDays = days
		}
	}

	batchSizeStr := os.Getenv("CLEANUP_BATCH_SIZE")
	batchSize := 100
	if batchSizeStr != "" {
		if size, err := strconv.Atoi(batchSizeStr); err == nil {
			batchSize = size
		}
	}

	return &CleanupConfig{
		RetentionPeriod: time.Duration(retentionDays*24) * time.Hour,
		BatchSize:       batchSize,
		DryRun:          os.Getenv("CLEANUP_DRY_RUN") == "true",
		ScheduleEnabled: os.Getenv("CLEANUP_SCHEDULE_ENABLED") == "true",
	}, nil
}

func (c *CleanupConfig) GetRetentionPeriod() time.Duration {
	return c.RetentionPeriod
}

func (c *CleanupConfig) GetBatchSize() int {
	return c.BatchSize
}

func (c *CleanupConfig) IsDryRun() bool {
	return c.DryRun
}

func (c *CleanupConfig) IsScheduleEnabled() bool {
	return c.ScheduleEnabled
}
