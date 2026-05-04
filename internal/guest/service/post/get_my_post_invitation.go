package post

import (
	"context"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
)

func (s *service) GetMyPostInvitations(ctx context.Context, customerID string) ([]entity.PostCollaborator, error) {
	return s.postRepo.GetPendingPostInvitations(ctx, customerID)
}
