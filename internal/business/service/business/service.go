package business

import (
	"context"

	"github.com/ua-academy-projects/share-bite/internal/business/entity"
	"github.com/ua-academy-projects/share-bite/pkg/database"
)

type businessRepository interface {
	UpdatePost(ctx context.Context, postID int64, orgID int, content string) (*entity.Post, error)
	DeletePost(ctx context.Context, id int64, orgID int) error
	GetOrgIDByUserID(ctx context.Context, userID string) (int, error)
	GetById(ctx context.Context, id int) (*entity.OrgUnit, error)
	ListByParentID(ctx context.Context, parentID, offset, limit int) ([]entity.OrgUnit, error)
	GetPostPhotos(ctx context.Context, postID int64) ([]string, error)
	CheckOwnership(ctx context.Context, userID string, unitID int) error
	CreatePost(ctx context.Context, userID string, unitID int, description string) (*entity.Post, error)
	InsertPostImages(ctx context.Context, postID int, URLs []string) error
	GetPosts(ctx context.Context, limit, offset int) ([]entity.Post, error)
	GetPostByID(ctx context.Context, postID int64) (*entity.Post, error)
}

type service struct {
	businessRepo businessRepository
	txManager    database.TxManager
}

func New(businessRepo businessRepository, txManager database.TxManager) *service {
	return &service{
		businessRepo: businessRepo,
		txManager:    txManager,
	}
}
