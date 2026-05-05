package post

import "context"

func (s *service) CleanupExpiredPosts(ctx context.Context) error {
	return s.postRepo.DeleteExpiredDraftPosts(ctx)
}
