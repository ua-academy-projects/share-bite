package business

import (
	"context"

	"github.com/ua-academy-projects/share-bite/internal/business/entity"
)

type businessRepository interface {
	UpdatePost(ctx context.Context, postID int64, orgID int, content string) (*entity.Post, error)
	DeletePost(ctx context.Context, id int64, orgID int) error
	GetOrgIDByUserID(ctx context.Context, userID int64) (int, error)
	GetById(ctx context.Context, id int) (*entity.OrgUnit, error)
	ListByParentID(ctx context.Context, parentID, offset, limit int) ([]entity.OrgUnit, error)
}

type service struct {
	businessRepo businessRepository
}

func New(businessRepo businessRepository) *service {
	return &service{
		businessRepo: businessRepo,
	}
}
