package business

import (
	"context"
	"fmt"

	"github.com/ua-academy-projects/share-bite/internal/business/entity"
	"github.com/ua-academy-projects/share-bite/pkg/database/pagination"
)

func (s *service) UpdatePost(ctx context.Context, postID int64, userID string, content string) (*entity.PostWithPhotos, error) {
	post, err := s.businessRepo.GetPostByID(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("get post: %w", err)
	}

	err = s.businessRepo.CheckOwnership(ctx, userID, post.OrgID)
	if err != nil {
		return nil, fmt.Errorf("check ownership: %w", err)
	}

	post, err = s.businessRepo.UpdatePost(ctx, postID, post.OrgID, content)
	if err != nil {
		return nil, fmt.Errorf("update post: %w", err)
	}

	photos, err := s.businessRepo.GetPostPhotos(ctx, post.ID)
	if err != nil {
		return nil, fmt.Errorf("get photos: %w", err)
	}

	return &entity.PostWithPhotos{
		ID:        post.ID,
		OrgID:     post.OrgID,
		Content:   post.Content,
		CreatedAt: post.CreatedAt,
		Images:    photos,
	}, nil
}

func (s *service) DeletePost(ctx context.Context, postID int64, userID string) error {
	post, err := s.businessRepo.GetPostByID(ctx, postID)
	if err != nil {
		return fmt.Errorf("get post: %w", err)
	}

	err = s.businessRepo.CheckOwnership(ctx, userID, post.OrgID)
	if err != nil {
		return fmt.Errorf("check ownership: %w", err)
	}

	return s.businessRepo.DeletePost(ctx, postID, post.OrgID)
}

func (s *service) CheckOwnership(ctx context.Context, userID string, unitID int) error {
	return s.businessRepo.CheckOwnership(ctx, userID, unitID)
}

func (s *service) CreatePost(ctx context.Context, userID string, unitID int, description string, URLs []string) (*entity.Post, error) {
	var post *entity.Post
	err := s.txManager.ReadCommited(ctx, func(ctxTx context.Context) error {
		var err error

		post, err = s.businessRepo.CreatePost(ctxTx, userID, unitID, description)
		if err != nil {
			return fmt.Errorf("create post: %w", err)
		}
		err = s.businessRepo.InsertPostImages(ctxTx, post.ID, URLs)
		if err != nil {
			return fmt.Errorf("insert images: %w", err)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("business service: %w", err)
	}
	return post, nil
}

func (s *service) GetPosts(ctx context.Context, skip, limit int) (pagination.Result[entity.PostWithPhotos], error) {

	const maxLimit = 100

	if skip < 0 {
		skip = 0
	}

	if limit < 1 {
		limit = 10
	}

	if limit > maxLimit {
		limit = maxLimit
	}

	posts, err := s.businessRepo.GetPosts(ctx, limit, skip)
	if err != nil {
		return pagination.Result[entity.PostWithPhotos]{}, fmt.Errorf("get posts: %w", err)
	}

	orgCache := make(map[int]*entity.OrgUnit)
	var items []entity.PostWithPhotos

	for _, post := range posts.Items {

		photos, err := s.businessRepo.GetPostPhotos(ctx, post.ID)
		if err != nil {
			return pagination.Result[entity.PostWithPhotos]{}, fmt.Errorf("get photos: %w", err)
		}

		org, ok := orgCache[post.OrgID]
		if !ok {
			org, err = s.businessRepo.GetById(ctx, post.OrgID)
			if err != nil {
				return pagination.Result[entity.PostWithPhotos]{}, fmt.Errorf("get org: %w", err)
			}
			orgCache[post.OrgID] = org
		}

		items = append(items, entity.PostWithPhotos{
			ID:          post.ID,
			OrgID:       post.OrgID,
			Content:     post.Content,
			CreatedAt:   post.CreatedAt,
			Images:      photos,
			OrgName:     org.Name,
			ProfileType: org.ProfileType,
		})
	}

	return pagination.Result[entity.PostWithPhotos]{
		Items: items,
		Total: posts.Total,
	}, nil
}
