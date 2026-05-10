package service

import (
	"context"
	"fmt"
	"strings"

	notificationentity "github.com/ua-academy-projects/share-bite/internal/notification/entity"
	notificationrepo "github.com/ua-academy-projects/share-bite/internal/notification/repository"
	"github.com/ua-academy-projects/share-bite/pkg/logger"
	"github.com/ua-academy-projects/share-bite/pkg/notification"
)

type Publisher interface {
	Publish(ctx context.Context, ch string, msg notification.Message) error
}

type Service struct {
	repo      notificationrepo.NotificationRepository
	publisher Publisher
}

func New(repo notificationrepo.NotificationRepository, publisher Publisher) *Service {
	return &Service{
		repo:      repo,
		publisher: publisher,
	}
}

func (s *Service) ProcessMessage(ctx context.Context, msg notification.Message) error {
	if s == nil {
		return fmt.Errorf("notification service is nil")
	}
	if msg.EventID == "" {
		return fmt.Errorf("event_id is required")
	}
	if msg.EventType == "" {
		return fmt.Errorf("event_type is required")
	}
	if msg.RecipientID == "" {
		return fmt.Errorf("recipient_id is required")
	}

	// We no longer call enrich here. The service is now generic.
	// It expects the sender to provide all necessary metadata (actor_name, actor_avatar, etc.)

	inserted, err := s.repo.Save(ctx, notificationentity.FromMessage(msg))
	if err != nil {
		return err
	}
	if !inserted {
		logger.DebugKV(ctx, "duplicate notification skipped", "notification_id", msg.EventID, "recipient_id", msg.RecipientID, "event_type", msg.EventType)
		return nil
	}

	if s.publisher != nil {
		if err := s.publisher.Publish(ctx, msg.RecipientID, msg); err != nil {
			return fmt.Errorf("publish notification to stream: %w", err)
		}
	}

	logger.InfoKV(ctx, "notification processed", "notification_id", msg.EventID, "recipient_id", msg.RecipientID, "event_type", msg.EventType)
	return nil
}

func mergeMetadata(existing, enriched map[string]any) map[string]any {
	if len(existing) == 0 && len(enriched) == 0 {
		return nil
	}

	merged := make(map[string]any, len(existing)+len(enriched))
	for key, value := range existing {
		merged[key] = value
	}
	for key, value := range enriched {
		if current, ok := existing[key]; ok && hasMeaningfulMetadataValue(current) {
			continue
		}
		merged[key] = value
	}

	return merged
}

func hasMeaningfulMetadataValue(value any) bool {
	switch typed := value.(type) {
	case nil:
		return false
	case string:
		return strings.TrimSpace(typed) != ""
	case []string:
		return len(typed) > 0
	default:
		return true
	}
}

func (s *Service) GetHistory(ctx context.Context, recipientID string, limit, offset int) ([]notificationentity.NotificationDTO, error) {
	items, err := s.repo.GetHistory(ctx, recipientID, limit, offset)
	if err != nil {
		return nil, err
	}

	result := make([]notificationentity.NotificationDTO, 0, len(items))
	for _, item := range items {
		result = append(result, item.ToDTO())
	}

	return result, nil
}

func (s *Service) MarkAsRead(ctx context.Context, recipientID string, notificationIDs []string) error {
	return s.repo.MarkAsRead(ctx, recipientID, notificationIDs)
}
