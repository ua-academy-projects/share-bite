package business

import (
	"context"

	"github.com/ua-academy-projects/share-bite/internal/business/entity"
	"github.com/ua-academy-projects/share-bite/pkg/database"
	"github.com/ua-academy-projects/share-bite/pkg/database/pagination"
)

type businessRepository interface {
	UpdatePost(ctx context.Context, postID int64, orgID int, content string) (*entity.Post, error)
	DeletePost(ctx context.Context, id int64, orgID int) error
	GetOrgIDByUserID(ctx context.Context, userID string) (int, error)
	GetById(ctx context.Context, id int) (*entity.OrgUnit, error)
	ListByParentID(ctx context.Context, parentID, offset, limit int) (pagination.Result[entity.OrgUnit], error)
	GetPostPhotos(ctx context.Context, postID int64) ([]string, error)
	CheckOwnership(ctx context.Context, userID string, unitID int) error
	CreatePost(ctx context.Context, userID string, unitID int, description string) (*entity.Post, error)
	InsertPostImages(ctx context.Context, postID int64, URLs []string) error
	GetPosts(ctx context.Context, skip, limit int) (pagination.Result[entity.Post], error)
	GetPostByID(ctx context.Context, postID int64) (*entity.Post, error)

	CreateLike(ctx context.Context, postID int64, customerID string) (*entity.Like, error)
	DeleteLike(ctx context.Context, postID int64, customerID string) error
	CheckUserLiked(ctx context.Context, postID int64, customerID string) (bool, error)
	CountLikesByPost(ctx context.Context, postID int64) (int, error)
	GetLikesByPost(ctx context.Context, postID int64, limit, offset int) ([]entity.Like, error)

	CreateComment(ctx context.Context, postID int64, authorID, content string) (*entity.Comment, error)
	GetCommentByID(ctx context.Context, commentID int64) (*entity.Comment, error)
	UpdateComment(ctx context.Context, commentID int64, content string) (*entity.Comment, error)
	DeleteComment(ctx context.Context, commentID int64) error
	ListCommentsWithAuthorsByPost(ctx context.Context, postID int64, limit, offset int) ([]entity.CommentWithAuthor, error)
	CountCommentsByPost(ctx context.Context, postID int64) (int, error)
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
