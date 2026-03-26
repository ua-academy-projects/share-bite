package business

import (
	"context"

	"github.com/ua-academy-projects/share-bite/internal/business/entity"
)

type businessRepository interface {
	UpdatePost(ctx context.Context, post *entity.Post) error
	DeletePost(ctx context.Context, id int64, orgID int) error
}

type service struct {
	businessRepo businessRepository
}

func New(businessRepo businessRepository) *service {
	return &service{
		businessRepo: businessRepo,
	}
}

func (s *service) UpdatePost(ctx context.Context, postID int64, orgID int, content string) error {
	post := &entity.Post{
		ID:      postID,
		OrgID:   orgID,
		Content: content,
	}

	return s.businessRepo.UpdatePost(ctx, post)
}

func (s *service) DeletePost(ctx context.Context, postID int64, orgID int) error {
	return s.businessRepo.DeletePost(ctx, postID, orgID)
}
