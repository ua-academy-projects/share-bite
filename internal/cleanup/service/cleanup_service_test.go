package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ua-academy-projects/share-bite/internal/cleanup/entity"
)

type MockCleanupRepository struct {
	mock.Mock
}

func (m *MockCleanupRepository) ExpireOldPosts(ctx context.Context, olderThan time.Time, batchSize int, dryRun bool) (int64, error) {
	args := m.Called(ctx, olderThan, batchSize, dryRun)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockCleanupRepository) CountExpiredPosts(ctx context.Context, olderThan time.Time) (int64, error) {
	args := m.Called(ctx, olderThan)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockCleanupRepository) GetExpiredPosts(ctx context.Context, olderThan time.Time, limit int, offset int) ([]*entity.ExpiredPost, error) {
	args := m.Called(ctx, olderThan, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.ExpiredPost), args.Error(1)
}

func (m *MockCleanupRepository) DeletePostsByID(ctx context.Context, postIDs []int64) (int64, error) {
	args := m.Called(ctx, postIDs)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockCleanupRepository) CountPasswordResetTokens(ctx context.Context, expiredBefore time.Time) (int64, error) {
	args := m.Called(ctx, expiredBefore)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockCleanupRepository) DeleteExpiredPasswordResetTokens(ctx context.Context, expiredBefore time.Time, batchSize int, dryRun bool) (int64, error) {
	args := m.Called(ctx, expiredBefore, batchSize, dryRun)
	return args.Get(0).(int64), args.Error(1)
}

func TestExpireOldPostsDryRun(t *testing.T) {
	mockRepo := new(MockCleanupRepository)
	mockRepo.On("CountExpiredPosts", mock.Anything, mock.Anything).Return(int64(150), nil)

	service := NewCleanupService(mockRepo)
	ctx := context.Background()
	retention := 30 * 24 * time.Hour

	result, err := service.ExpireOldPosts(ctx, retention, 100, true)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "expire-old-posts", result.Name)
	assert.True(t, result.DryRun)
	assert.Equal(t, int64(150), result.RecordsFound)
	assert.Equal(t, int64(150), result.RecordsDeleted)
	assert.Empty(t, result.Errors)

	mockRepo.AssertCalled(t, "CountExpiredPosts", mock.Anything, mock.Anything)
}

func TestExpireOldPostsProduction(t *testing.T) {
	mockRepo := new(MockCleanupRepository)
	mockRepo.On("CountExpiredPosts", mock.Anything, mock.Anything).Return(int64(250), nil)
	mockRepo.On("ExpireOldPosts", mock.Anything, mock.Anything, mock.Anything, false).Return(int64(250), nil)

	service := NewCleanupService(mockRepo)
	ctx := context.Background()
	retention := 30 * 24 * time.Hour

	result, err := service.ExpireOldPosts(ctx, retention, 100, false)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "expire-old-posts", result.Name)
	assert.False(t, result.DryRun)
	assert.Equal(t, int64(250), result.RecordsFound)
	assert.Equal(t, int64(250), result.RecordsDeleted)
	assert.Empty(t, result.Errors)

	mockRepo.AssertCalled(t, "CountExpiredPosts", mock.Anything, mock.Anything)
	mockRepo.AssertCalled(t, "ExpireOldPosts", mock.Anything, mock.Anything, 100, false)
}

func TestExpireOldPostsNoRecords(t *testing.T) {
	// Setup
	mockRepo := new(MockCleanupRepository)
	mockRepo.On("CountExpiredPosts", mock.Anything, mock.Anything).Return(int64(0), nil)

	service := NewCleanupService(mockRepo)
	ctx := context.Background()
	retention := 30 * 24 * time.Hour

	// Execute
	result, err := service.ExpireOldPosts(ctx, retention, 100, false)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(0), result.RecordsFound)
	assert.Equal(t, int64(0), result.RecordsDeleted)
	assert.Empty(t, result.Errors)
}

func TestCleanExpiredPasswordResetTokensDryRun(t *testing.T) {
	// Setup
	mockRepo := new(MockCleanupRepository)
	mockRepo.On("CountPasswordResetTokens", mock.Anything, mock.Anything).Return(int64(50), nil)

	service := NewCleanupService(mockRepo)
	ctx := context.Background()
	retention := 24 * time.Hour

	// Execute
	result, err := service.CleanExpiredPasswordResetTokens(ctx, retention, 100, true)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "cleanup-expired-tokens", result.Name)
	assert.True(t, result.DryRun)
	assert.Equal(t, int64(50), result.RecordsFound)
	assert.Equal(t, int64(50), result.RecordsDeleted)
	assert.Empty(t, result.Errors)

	mockRepo.AssertCalled(t, "CountPasswordResetTokens", mock.Anything, mock.Anything)
}

func TestCleanExpiredPasswordResetTokensProduction(t *testing.T) {
	// Setup
	mockRepo := new(MockCleanupRepository)
	mockRepo.On("CountPasswordResetTokens", mock.Anything, mock.Anything).Return(int64(75), nil)
	mockRepo.On("DeleteExpiredPasswordResetTokens", mock.Anything, mock.Anything, mock.Anything, false).Return(int64(75), nil)

	service := NewCleanupService(mockRepo)
	ctx := context.Background()
	retention := 24 * time.Hour

	// Execute
	result, err := service.CleanExpiredPasswordResetTokens(ctx, retention, 100, false)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "cleanup-expired-tokens", result.Name)
	assert.False(t, result.DryRun)
	assert.Equal(t, int64(75), result.RecordsFound)
	assert.Equal(t, int64(75), result.RecordsDeleted)
	assert.Empty(t, result.Errors)

	mockRepo.AssertCalled(t, "CountPasswordResetTokens", mock.Anything, mock.Anything)
	mockRepo.AssertCalled(t, "DeleteExpiredPasswordResetTokens", mock.Anything, mock.Anything, 100, false)
}

func TestRunAllCleanups(t *testing.T) {
	// Setup
	mockRepo := new(MockCleanupRepository)
	mockRepo.On("CountExpiredPosts", mock.Anything, mock.Anything).Return(int64(100), nil)
	mockRepo.On("ExpireOldPosts", mock.Anything, mock.Anything, 100, false).Return(int64(100), nil)
	mockRepo.On("CountPasswordResetTokens", mock.Anything, mock.Anything).Return(int64(50), nil)
	mockRepo.On("DeleteExpiredPasswordResetTokens", mock.Anything, mock.Anything, 100, false).Return(int64(50), nil)

	service := NewCleanupService(mockRepo)
	ctx := context.Background()
	retention := 30 * 24 * time.Hour

	// Execute
	results, err := service.RunAllCleanups(ctx, retention, 100, false)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, results)
	assert.Equal(t, 2, len(results))

	assert.Equal(t, "expire-old-posts", results[0].Name)
	assert.Equal(t, int64(100), results[0].RecordsFound)
	assert.Equal(t, int64(100), results[0].RecordsDeleted)

	assert.Equal(t, "cleanup-expired-tokens", results[1].Name)
	assert.Equal(t, int64(50), results[1].RecordsFound)
	assert.Equal(t, int64(50), results[1].RecordsDeleted)
}

func TestCleanupResultDuration(t *testing.T) {
	now := time.Now()
	result := &entity.CleanupResult{
		Name:        "test-job",
		StartedAt:   now,
		CompletedAt: now.Add(5 * time.Second),
	}

	duration := result.Duration()
	assert.Equal(t, 5*time.Second, duration)
}

func TestCleanupResultWithErrors(t *testing.T) {
	// Setup
	mockRepo := new(MockCleanupRepository)
	mockRepo.On("CountExpiredPosts", mock.Anything, mock.Anything).Return(int64(0), assert.AnError)

	service := NewCleanupService(mockRepo)
	ctx := context.Background()
	retention := 30 * 24 * time.Hour

	// Execute
	result, err := service.ExpireOldPosts(ctx, retention, 100, false)

	// Assert
	assert.Error(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.Errors)
}
