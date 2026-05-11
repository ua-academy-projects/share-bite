package entity

import "time"

type CleanupResult struct {
	Name           string
	RecordsFound   int64
	RecordsDeleted int64
	DryRun         bool
	StartedAt      time.Time
	CompletedAt    time.Time
	Errors         []string
}

func (r *CleanupResult) Duration() time.Duration {
	return r.CompletedAt.Sub(r.StartedAt)
}

type CleanupJob struct {
	Name         string
	Description  string
	RetentionAge time.Duration
	BatchSize    int
	DryRun       bool
}

type ExpiredPost struct {
	ID         int64
	CustomerID string
	CreatedAt  time.Time
}

type ExpiredPostCount struct {
	Count int64
}
