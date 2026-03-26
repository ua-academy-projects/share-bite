package business

import (
	"context"
)

type businessRepository interface {
	UpdatePost(ctx context.Context, postID int64, orgID int, content string) error
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

	return s.businessRepo.UpdatePost(ctx, postID, orgID, content)
}

func (s *service) DeletePost(ctx context.Context, postID int64, userID int64) error {
	orgID, err := s.businessRepo.GetOrgIDByUserID(ctx, userID)
	if err != nil {
		return err
	}

	return s.businessRepo.DeletePost(ctx, postID, orgID)
}
