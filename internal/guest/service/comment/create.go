package comment

import (
	"context"
	"strconv"

	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	"github.com/ua-academy-projects/share-bite/pkg/errwrap"
)

func (s *service) Create(ctx context.Context, in entity.CreateCommentInput) (entity.Comment, error) {
	postIDStr := strconv.FormatInt(in.PostID, 10)

	_, err := s.postSvc.Get(ctx, postIDStr)
	if err != nil {
		return entity.Comment{}, errwrap.Wrap("check post existence in comment service", err)
	}

	comment, err := s.commentRepo.Create(ctx, in)
	if err != nil {
		return entity.Comment{}, errwrap.Wrap("create comment in repo", err)
	}

	return comment, nil
}
