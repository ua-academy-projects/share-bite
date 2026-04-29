package worker

import (
	"context"

	"github.com/robfig/cron/v3"
	"github.com/ua-academy-projects/share-bite/pkg/logger"
)

type TokenRepository interface {
	DeleteExpiredTokens(ctx context.Context) error
}

type Manager struct {
	scheduler *cron.Cron
	repo      TokenRepository
}

func NewManager(repo TokenRepository) *Manager {
	return &Manager{
		scheduler: cron.New(),
		repo:      repo,
	}
}

func (m *Manager) Start(ctx context.Context) {
	_, err := m.scheduler.AddFunc("0 3 * * *", func() {
		logger.Info(ctx, "CRON: Starting expired tokens cleanup")

		if err := m.repo.DeleteExpiredTokens(ctx); err != nil {
			logger.Error(ctx, "CRON: Failed to delete expired tokens", err)
			return
		}

		logger.Info(ctx, "CRON: Expired tokens cleanup finished successfully")
	})

	if err != nil {
		logger.Fatal(ctx, "CRON: Failed to schedule job", err)
	}

	m.scheduler.Start()
	logger.Info(ctx, "CRON: Scheduler started")
}

func (m *Manager) Stop() {
	m.scheduler.Stop()
}
