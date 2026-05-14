package business

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/google/uuid"
	"github.com/ua-academy-projects/share-bite/internal/business/entity"
	biserr "github.com/ua-academy-projects/share-bite/internal/business/error"
	"github.com/ua-academy-projects/share-bite/internal/storage/key"
	"github.com/ua-academy-projects/share-bite/internal/storage/mediatype"
	"github.com/ua-academy-projects/share-bite/pkg/database/pagination"
)

const (
	fileHeaderSize = 512
)

var (
	postImageValidator = mediatype.NewValidator(mediatype.DefaultMaxImageSizeBytes, "image/jpeg", "image/png", "image/webp", "image/gif")
)

func (s *service) CreatePost(ctx context.Context, userID string, unitID int, description string, images []*multipart.FileHeader) (*entity.PostWithPhotos, error) {
	const op = "service.post.CreatePost"
	for _, file := range images {
		openedFile, err := file.Open()
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		buffer := make([]byte, fileHeaderSize)
		_, err = openedFile.Read(buffer)
		openedFile.Close()

		if err != nil && err != io.EOF {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		contentType := http.DetectContentType(buffer)
		if err := postImageValidator.Validate(contentType, file.Size); err != nil {
			if errors.Is(err, mediatype.ErrUnsupportedType) {
				return nil, biserr.WrongFileExtErr
			}

			if errors.Is(err, mediatype.ErrFileTooLarge) {
				return nil, biserr.FileToLargeErr
			}

			return nil, err
		}
	}

	var photoURLs []string
	var uploadedKeys []string
	uploadSessionID := uuid.NewString()

	cleanupS3 := func() {
		for _, key := range uploadedKeys {
			s.storage.Delete(ctx, key)
		}
	}

	isSuccess := false
	defer func() {
		if !isSuccess {
			cleanupS3()
		}
	}()

	for _, file := range images {
		openedFile, err := file.Open()
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		buffer := make([]byte, fileHeaderSize)
		_, err = openedFile.Read(buffer)
		if err != nil && err != io.EOF {
			openedFile.Close()
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		contentType := http.DetectContentType(buffer)
		ext, ok := mediatype.ExtFromContentType(contentType)
		if !ok {
			openedFile.Close()
			return nil, biserr.WrongFileExtErr
		}

		seeker, ok := openedFile.(io.Seeker)
		if !ok {
			openedFile.Close()
			return nil, fmt.Errorf("%s: uploaded file is not seekable", op)
		}
		if _, err := seeker.Seek(0, io.SeekStart); err != nil {
			openedFile.Close()
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		objectKey := key.BusinessPostImageKey(unitID, uploadSessionID, uuid.NewString(), ext)
		err = s.storage.Upload(ctx, objectKey, contentType, openedFile)

		openedFile.Close()

		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		uploadedKeys = append(uploadedKeys, objectKey)

		imageURL := s.storage.BuildURL(objectKey)
		photoURLs = append(photoURLs, imageURL)
	}

	var post *entity.Post

	err := s.txManager.ReadCommitted(ctx, func(ctxTx context.Context) error {
		var err error

		post, err = s.businessRepo.CreatePost(ctxTx, userID, unitID, description)
		if err != nil {
			return fmt.Errorf("create post: %w", err)
		}
		err = s.businessRepo.InsertPostImages(ctxTx, post.ID, uploadedKeys)
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

	isSuccess = true
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

func (s *service) UpdatePost(ctx context.Context, postID int64, userID string, content string) (*entity.PostWithPhotos, error) {
	const op = "service.post.UpdatePost"
	post, err := s.businessRepo.GetPostByID(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = s.businessRepo.CheckOwnership(ctx, userID, post.OrgID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	post, err = s.businessRepo.UpdatePost(ctx, postID, post.OrgID, content)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	photos, err := s.businessRepo.GetPostPhotos(ctx, post.ID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
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
	const op = "service.post.DeletePost"
	post, err := s.businessRepo.GetPostByID(ctx, postID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = s.businessRepo.CheckOwnership(ctx, userID, post.OrgID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = s.businessRepo.DeletePost(ctx, postID, post.OrgID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *service) CheckOwnership(ctx context.Context, userID string, unitID int) error {
	const op = "service.post.CheckOwnership"

	err := s.businessRepo.CheckOwnership(ctx, userID, unitID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *service) GetPosts(ctx context.Context, skip, limit int, orgIDs []int) (pagination.Result[entity.PostWithPhotos], error) {
	const op = "service.post.GetPosts"
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

	posts, err := s.businessRepo.GetPosts(ctx, limit, skip, orgIDs)
	if err != nil {
		return pagination.Result[entity.PostWithPhotos]{}, fmt.Errorf("%s: %w", op, err)
	}

	orgCache := make(map[int]*entity.OrgUnit)
	var items []entity.PostWithPhotos

	for _, post := range posts.Items {

		photos, err := s.businessRepo.GetPostPhotos(ctx, post.ID)
		if err != nil {
			return pagination.Result[entity.PostWithPhotos]{}, fmt.Errorf("%s: %w", op, err)
		}

		org, ok := orgCache[post.OrgID]
		if !ok {
			org, err = s.businessRepo.GetById(ctx, post.OrgID)
			if err != nil {
				return pagination.Result[entity.PostWithPhotos]{}, fmt.Errorf("%s: %w", op, err)
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
