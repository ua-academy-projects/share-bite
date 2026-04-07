package business

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/ua-academy-projects/share-bite/internal/business/entity"
	biserr "github.com/ua-academy-projects/share-bite/internal/business/error"
	"github.com/ua-academy-projects/share-bite/pkg/database/pagination"
)

const (
	maxImageSize   = 5 * 1024 * 1024
	fileHeaderSize = 512
)

func (s *service) CreatePost(ctx context.Context, userID string, unitID int, description string, images []*multipart.FileHeader) (*entity.PostWithPhotos, error) {
	for _, file := range images {
		if file.Size > maxImageSize {
			return nil, biserr.FileToLargeErr
		}

		openedFile, err := file.Open()
		if err != nil {
			return nil, err
		}

		buffer := make([]byte, fileHeaderSize)
		_, err = openedFile.Read(buffer)
		openedFile.Close()

		if err != nil && err != io.EOF {
			return nil, err
		}
		contentType := http.DetectContentType(buffer)
		if !isAllowedImageType(contentType) {
			return nil, biserr.WrongFileExtErr
		}
	}

	var photoURLs []string
	var uploadedKeys []string

	cleanupS3 := func() {
		for _, key := range uploadedKeys {
			s.storage.Delete(ctx, key)
		}
	}

	for _, file := range images {
		openedFile, err := file.Open()
		if err != nil {
			cleanupS3()
			return nil, err
		}

		buffer := make([]byte, fileHeaderSize)
		_, err = openedFile.Read(buffer)

		fileExt := filepath.Ext(file.Filename)
		contentType := http.DetectContentType(buffer)

		seeker, _ := openedFile.(io.Seeker)
		seeker.Seek(0, io.SeekStart)

		objectKey := fmt.Sprintf("posts/%d/%s%s", unitID, uuid.New().String(), fileExt)

		key, err := s.storage.Upload(ctx, objectKey, contentType, openedFile)

		openedFile.Close()

		if err != nil {
			cleanupS3()
			return nil, err
		}

		uploadedKeys = append(uploadedKeys, key)
		imageURL := s.storage.BuildURL(key)
		photoURLs = append(photoURLs, imageURL)
	}

	var post *entity.Post

	err := s.txManager.ReadCommited(ctx, func(ctxTx context.Context) error {
		var err error

		post, err = s.businessRepo.CreatePost(ctxTx, userID, unitID, description)
		if err != nil {
			return fmt.Errorf("create post: %w", err)
		}
		err = s.businessRepo.InsertPostImages(ctxTx, post.ID, photoURLs)
		if err != nil {
			return fmt.Errorf("insert images: %w", err)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("business service: %w", err)
	}

	org, err := s.businessRepo.GetById(ctx, post.OrgID)
	if err != nil {
		return nil, fmt.Errorf("get org: %w", err)
	}

	return &entity.PostWithPhotos{
		ID:          post.ID,
		OrgID:       post.OrgID,
		Content:     post.Content,
		CreatedAt:   post.CreatedAt,
		Images:      photoURLs,
		OrgName:     org.Name,
		ProfileType: org.ProfileType,
	}, nil
}

func isAllowedImageType(contentType string) bool {
	switch contentType {
	case "image/jpeg", "image/jpg", "image/png", "image/webp", "image/gif":
		return true
	default:
		return false
	}
}

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
