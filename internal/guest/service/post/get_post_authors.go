package post

import "context"

func (s *service) GetPostAuthors(ctx context.Context, postID string) ([]string, error) {
	collaborators, err := s.postRepo.GetAcceptedPostCollaborators(ctx, postID)
	if err != nil {
		return nil, err
	}

	authorID, err := s.postRepo.GetAuthorUserID(ctx, postID)
	if err != nil {
		return nil, err
	}

	return append([]string{authorID}, collaborators...), nil
}
