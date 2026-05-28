package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCleanupResultDuration(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name     string
		result   *CleanupResult
		expected time.Duration
	}{
		{
			name: "5 second duration",
			result: &CleanupResult{
				StartedAt:   now,
				CompletedAt: now.Add(5 * time.Second),
			},
			expected: 5 * time.Second,
		},
		{
			name: "1 minute duration",
			result: &CleanupResult{
				StartedAt:   now,
				CompletedAt: now.Add(1 * time.Minute),
			},
			expected: 1 * time.Minute,
		},
		{
			name: "0 duration",
			result: &CleanupResult{
				StartedAt:   now,
				CompletedAt: now,
			},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.result.Duration())
		})
	}
}

func TestCleanupResult(t *testing.T) {
	result := &CleanupResult{
		Name:           "test-job",
		RecordsFound:   100,
		RecordsDeleted: 100,
		DryRun:         false,
		StartedAt:      time.Now(),
		CompletedAt:    time.Now().Add(10 * time.Second),
		Errors:         []string{},
	}

	assert.Equal(t, "test-job", result.Name)
	assert.Equal(t, int64(100), result.RecordsFound)
	assert.Equal(t, int64(100), result.RecordsDeleted)
	assert.False(t, result.DryRun)
	assert.Empty(t, result.Errors)
	assert.Equal(t, 10*time.Second, result.Duration())
}

func TestCleanupResultWithErrors(t *testing.T) {
	errors := []string{"error 1", "error 2"}
	result := &CleanupResult{
		Name:        "test-job",
		DryRun:      false,
		Errors:      errors,
		StartedAt:   time.Now(),
		CompletedAt: time.Now(),
	}

	assert.Equal(t, 2, len(result.Errors))
	assert.Equal(t, "error 1", result.Errors[0])
	assert.Equal(t, "error 2", result.Errors[1])
}

func TestCleanupJob(t *testing.T) {
	job := &CleanupJob{
		Name:         "expire-posts",
		Description:  "Expire old posts",
		RetentionAge: 30 * 24 * time.Hour,
		BatchSize:    100,
		DryRun:       false,
	}

	assert.Equal(t, "expire-posts", job.Name)
	assert.Equal(t, "Expire old posts", job.Description)
	assert.Equal(t, 30*24*time.Hour, job.RetentionAge)
	assert.Equal(t, 100, job.BatchSize)
	assert.False(t, job.DryRun)
}

func TestExpiredPost(t *testing.T) {
	now := time.Now()
	post := &ExpiredPost{
		ID:         123,
		CustomerID: "cust-123",
		CreatedAt:  now,
	}

	assert.Equal(t, int64(123), post.ID)
	assert.Equal(t, "cust-123", post.CustomerID)
	assert.Equal(t, now, post.CreatedAt)
}

func TestExpiredPostCount(t *testing.T) {
	count := &ExpiredPostCount{
		Count: 50,
	}

	assert.Equal(t, int64(50), count.Count)
}
