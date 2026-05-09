package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	notificationentity "github.com/ua-academy-projects/share-bite/internal/notification/entity"
	notificationrepo "github.com/ua-academy-projects/share-bite/internal/notification/repository"
	"github.com/ua-academy-projects/share-bite/pkg/logger"
	"github.com/ua-academy-projects/share-bite/pkg/notification"
)

type CustomerProvider interface {
	GetByUserID(ctx context.Context, userID string) (entity.Customer, error)
}

type AvatarURLBuilder interface {
	BuildURL(key string) string
}

type Publisher interface {
	Publish(ctx context.Context, ch string, msg notification.Message) error
}

type Service struct {
	repo           notificationrepo.NotificationRepository
	customers      CustomerProvider
	publisher      Publisher
	avatarURLBuild AvatarURLBuilder
}

func New(repo notificationrepo.NotificationRepository, customers CustomerProvider, publisher Publisher, avatarURLBuild AvatarURLBuilder) *Service {
	return &Service{
		repo:           repo,
		customers:      customers,
		publisher:      publisher,
		avatarURLBuild: avatarURLBuild,
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

	metadata, err := s.enrich(ctx, msg)
	if err != nil {
		return err
	}
	msg.Metadata = metadata

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

func (s *Service) enrich(ctx context.Context, msg notification.Message) (map[string]any, error) {
	switch msg.EventType {
	case notification.PostLiked:
		return s.enrichPostLiked(ctx, msg)
	case notification.RegistrationConfirmed:
		return s.enrichRegistrationConfirmed(ctx, msg)
	default:
		return nil, fmt.Errorf("unknown notification event type: %s", msg.EventType)
	}
}

func (s *Service) enrichPostLiked(ctx context.Context, msg notification.Message) (map[string]any, error) {
	actor, err := s.getActor(ctx, msg.ActorID)
	if err != nil {
		logger.WarnKV(ctx, "post liked metadata fallback", "actor_id", msg.ActorID, "error", err)
		actor = actorProfile{name: "Share Bite"}
	}

	return map[string]any{
		"actor_name":     actor.name,
		"actor_avatar":   actor.avatar,
		"actor_username": actor.username,
	}, nil
}

func (s *Service) enrichRegistrationConfirmed(ctx context.Context, msg notification.Message) (map[string]any, error) {
	actor, err := s.getActor(ctx, msg.ActorID)
	if err != nil {
		logger.WarnKV(ctx, "registration confirmed metadata fallback", "actor_id", msg.ActorID, "error", err)
		actor = actorProfile{name: "Share Bite"}
	}

	return map[string]any{
		"actor_name":     actor.name,
		"actor_avatar":   actor.avatar,
		"actor_username": actor.username,
	}, nil
}

type actorProfile struct {
	name     string
	username string
	avatar   string
}

func (s *Service) getActor(ctx context.Context, userID string) (actorProfile, error) {
	if userID == "" {
		return actorProfile{name: "Share Bite"}, nil
	}
	if s.customers == nil {
		return actorProfile{name: "Share Bite"}, nil
	}

	customer, err := s.customers.GetByUserID(ctx, userID)
	if err != nil {
		return actorProfile{}, fmt.Errorf("get actor customer %s: %w", userID, err)
	}

	name := strings.TrimSpace(customer.FirstName + " " + customer.LastName)
	if name == "" {
		name = customer.UserName
	}
	if name == "" {
		name = "Share Bite"
	}

	avatar := ""
	if customer.AvatarObjectKey != nil && s.avatarURLBuild != nil {
		avatar = s.avatarURLBuild.BuildURL(*customer.AvatarObjectKey)
	}

	return actorProfile{
		name:     name,
		username: customer.UserName,
		avatar:   avatar,
	}, nil
}
