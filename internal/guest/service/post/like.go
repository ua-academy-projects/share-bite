package post

import (
	"context"

	"github.com/ua-academy-projects/share-bite/pkg/errwrap"
)

func (s *service) Like(ctx context.Context, postID string, customerID string) error {
	if err := s.postRepo.Like(ctx, postID, customerID); err != nil {
		return errwrap.Wrap("like post in repository", err)
	}
	return nil
}

func (s *service) Unlike(ctx context.Context, postID string, customerID string) error {
	if err := s.postRepo.Unlike(ctx, postID, customerID); err != nil {
		return errwrap.Wrap("unlike post in repository", err)
	}
	return nil
}
