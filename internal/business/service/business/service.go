package business

import (
	"context"

	"github.com/ua-academy-projects/share-bite/internal/business/entity"
	biserr "github.com/ua-academy-projects/share-bite/internal/business/error"
)

type businessRepository interface {
	GetPostByID(ctx context.Context, id int) (*entity.Post, error)
	UpdatePost(ctx context.Context, post *entity.Post) error
	DeletePost(ctx context.Context, id int) error
}

type service struct {
	businessRepo businessRepository
}

func New(businessRepo businessRepository) *service {
	return &service{
		businessRepo: businessRepo,
	}
}

func (s *service) UpdatePost(ctx context.Context, postID int, orgID int, content string) error {
	post, err := s.businessRepo.GetPostByID(ctx, postID)
	if err != nil {
		return err
	}

	if post.OrgID != orgID {
		return biserr.ErrForbidden
	}

	post.Content = content

	return s.businessRepo.UpdatePost(ctx, post)
}

func (s *service) DeletePost(ctx context.Context, postID int, orgID int) error {
	post, err := s.businessRepo.GetPostByID(ctx, postID)
	if err != nil {
		return err
	}

	if post.OrgID != orgID {
		return biserr.ErrForbidden
	}

	return s.businessRepo.DeletePost(ctx, postID)
}
