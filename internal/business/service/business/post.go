package business

import (
	"context"

	"github.com/ua-academy-projects/share-bite/internal/business/entity"
)

func (s *service) UpdatePost(ctx context.Context, postID int64, userID int64, content string) (*entity.PostWithPhotos, error) {
	orgID, err := s.businessRepo.GetOrgIDByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	post, err := s.businessRepo.UpdatePost(ctx, postID, orgID, content)
	if err != nil {
		return nil, err
	}

	photos, err := s.businessRepo.GetPostPhotos(ctx, post.ID)
	if err != nil {
		return nil, err
	}

	return &entity.PostWithPhotos{
		ID:        post.ID,
		OrgID:     post.OrgID,
		Content:   post.Content,
		CreatedAt: post.CreatedAt,
		Images:    photos,
	}, nil
}

func (s *service) DeletePost(ctx context.Context, postID int64, userID int64) error {
	orgID, err := s.businessRepo.GetOrgIDByUserID(ctx, userID)
	if err != nil {
		return err
	}

	return s.businessRepo.DeletePost(ctx, postID, orgID)
}
