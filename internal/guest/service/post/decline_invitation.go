package post

import "context"

func (s *service) DeclineInvitation(ctx context.Context, collaboratorID string, customerID string) error {
	_, err := s.postRepo.DeclinePostInvitation(ctx, collaboratorID, customerID)
	return err
}
