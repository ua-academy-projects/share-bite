package business

import (
	"context"

	"github.com/ua-academy-projects/share-bite/pkg/database"

	"github.com/ua-academy-projects/share-bite/internal/business/entity"
	"github.com/ua-academy-projects/share-bite/internal/storage"
)

type businessRepository interface {
	CheckOwnership(ctx context.Context, userID string, unitID int) error
	CreatePost(ctx context.Context, userID string, unitID int, description string) (*entity.Post, error)	
	InsertPostImages(ctx context.Context, postID int, URLs []string) error
}

type service struct {
	businessRepo businessRepository
	txManager    database.TxManager
	storage storage.ObjectStorage
}

func New(businessRepo businessRepository, txManager database.TxManager, st storage.ObjectStorage) *service {
	return &service{
		businessRepo: businessRepo,
		txManager:    txManager,
		storage: st,
	}
}

func (s *service) CheckOwnership(ctx context.Context, userID string, unitID int) error { //for handlers
	err := s.businessRepo.CheckOwnership(ctx, userID, unitID)
	if err != nil {
		return err
	}
	return nil
}

