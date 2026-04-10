package comment

import (
	"context"
	"fmt"
	"strconv"

	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
)

func (s *service) List(ctx context.Context, in entity.ListCommentsInput) (entity.ListCommentsOutput, error) {
	postIDStr := strconv.FormatInt(in.PostID, 10)
	if _, err := s.postSvc.Get(ctx, postIDStr, ""); err != nil {
		return entity.ListCommentsOutput{}, fmt.Errorf("check post existence for comments list: %w", err)
	}

	out, err := s.commentRepo.List(ctx, in)
	if err != nil {
		return entity.ListCommentsOutput{}, fmt.Errorf("get list of comments from repo: %w", err)
	}

	return out, nil
}
