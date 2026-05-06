package post

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/ua-academy-projects/share-bite/internal/guest/dto"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
)

func (s *service) Create(ctx context.Context, in dto.CreatePostInput) (entity.Post, error) {
	var post entity.Post

	err := s.txManager.ReadCommitted(ctx, func(txCtx context.Context) error {
		createdPost, err := s.createPost(txCtx, in)
		if err != nil {
			return err
		}

		post = createdPost
		return nil
	})

	if err != nil {
		return entity.Post{}, err
	}

	return post, nil
}

func generatePostImageKey(customerID, postID, ext string) string {
	return fmt.Sprintf("posts/%s/%s/%s.%s", customerID, postID, uuid.New().String(), ext)
}

func extensionFromContentType(contentType string) string {
	switch contentType {
	case "image/jpeg":
		return "jpg"
	case "image/png":
		return "png"
	default:
		return ""
	}
}

func UniqueStrings(input []string) []string {
	set := make(map[string]struct{}, len(input))
	res := make([]string, 0, len(input))
	for _, v := range input {
		if _, ok := set[v]; ok {
			continue
		}
		set[v] = struct{}{}
		res = append(res, v)
	}
	return res
}