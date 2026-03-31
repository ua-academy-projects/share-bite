package business

import (
	"context"

	"github.com/ua-academy-projects/share-bite/internal/business/entity"
)

func (s *service) UpdatePost(ctx context.Context, postID int64, userID int64, content string) (*entity.Post, error) {
	orgID, err := s.businessRepo.GetOrgIDByUserID(ctx, userID)
	if err != nil {
		return nil, err
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
