package comment

import (
	"context"
	"strconv"
	"fmt"

	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
)

func (s *service) Create(ctx context.Context, in entity.CreateCommentInput) (entity.Comment, error) {
	postIDStr := strconv.FormatInt(in.PostID, 10)

	_, err := s.postSvc.Get(ctx, postIDStr, "")
	if err != nil {
		return entity.Comment{}, fmt.Errorf("check post existence in comment service: %w", err)
	}

	comment, err := s.commentRepo.Create(ctx, in)
	if err != nil {
		return entity.Comment{}, fmt.Errorf("create comment in repo: %w", err)
	}

	return comment, nil
}
