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

	prefs, err := s.repo.GetPreferences(ctx, msg.RecipientID)
	if err != nil {
		return fmt.Errorf("get notification preferences: %w", err)
	}

	if enabled, ok := prefs[string(msg.EventType)]; ok && !enabled {
		logger.InfoKV(ctx, "notification skipped due to user preference", "notification_id", msg.EventID, "recipient_id", msg.RecipientID, "event_type", msg.EventType)
		return nil
	}

	inserted, err := s.repo.Save(ctx, notificationentity.FromMessage(msg))
	if err != nil {
		return err
	}
	if !inserted {
		logger.DebugKV(ctx, "duplicate notification skipped from db save", "notification_id", msg.EventID, "recipient_id", msg.RecipientID, "event_type", msg.EventType)
		if s.publisher == nil {
			return nil
		}
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

func (s *Service) GetPreferences(ctx context.Context, recipientID string) (map[string]bool, error) {
	dbPrefs, err := s.repo.GetPreferences(ctx, recipientID)
	if err != nil {
		return nil, err
	}

	result := map[string]bool{
		"post_liked":               true,
		"invitation_received":      true,
		"post_published":           true,
		"post_invitation_accepted": true,
		"business_verified":        true,
		"business_rejected":        true,
	}

	for k, v := range dbPrefs {
		if _, ok := result[k]; ok {
			result[k] = v
		}
	}

	return result, nil
}

func (s *Service) UpdatePreferences(ctx context.Context, recipientID string, prefs map[string]bool) error {
	validKeys := map[string]bool{
		"post_liked":               true,
		"invitation_received":      true,
		"post_published":           true,
		"post_invitation_accepted": true,
		"business_verified":        true,
		"business_rejected":        true,
	}

	for k := range prefs {
		if !validKeys[k] {
			return fmt.Errorf("unsupported preference key: %s", k)
		}
	}

	return s.repo.UpdatePreferences(ctx, recipientID, prefs)
}
