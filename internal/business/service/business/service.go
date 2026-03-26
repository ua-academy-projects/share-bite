package business

import (
	"context"

	"github.com/ua-academy-projects/share-bite/internal/business/entity"
)

type businessRepository interface {
	UpdatePost(ctx context.Context, post *entity.Post) error
	DeletePost(ctx context.Context, id int64, orgID int) error
	GetOrgIDByUserID(ctx context.Context, userID int64) (int, error)
}

type service struct {
	businessRepo businessRepository
}

func New(businessRepo businessRepository) *service {
	return &service{
		businessRepo: businessRepo,
	}
}

func (s *service) UpdatePost(ctx context.Context, postID int64, userID int64, content string) error {
	orgID, err := s.businessRepo.GetOrgIDByUserID(ctx, userID)
	if err != nil {
		return err
	}

	post := &entity.Post{
		ID:      postID,
		OrgID:   orgID,
		Content: content,
	}

	return s.businessRepo.UpdatePost(ctx, post)
}

func (s *service) DeletePost(ctx context.Context, postID int64, userID int64) error {
	orgID, err := s.businessRepo.GetOrgIDByUserID(ctx, userID)
	if err != nil {
		return err
	}

	return s.businessRepo.DeletePost(ctx, postID, orgID)
}
