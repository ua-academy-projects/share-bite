package post

import (
	"context"
	"fmt"
	"github.com/ua-academy-projects/share-bite/internal/guest/dto"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
)

func (s *service) List(ctx context.Context, in dto.ListPostsInput) (dto.ListPostsOutput, error) {
	out, err := s.postRepo.List(ctx, in)
	if err != nil {
		return dto.ListPostsOutput{}, fmt.Errorf("get list of posts from post repository: %w", err)
	}

	if len(out.Posts) == 0 {
		return out, nil
	}

	postIDs := make([]string, 0, len(out.Posts))
	for _, p := range out.Posts {
		postIDs = append(postIDs, p.ID)
	}

	mentionsMap, err := s.postRepo.ListMentionsByPostIDs(ctx, postIDs)
	if err != nil {
		return dto.ListPostsOutput{}, err
	}

	customerIDSet := make(map[string]struct{})

	for _, mentions := range mentionsMap {
		for _, m := range mentions {
			customerIDSet[m.CustomerID] = struct{}{}
		}
	}

	customerIDs := make([]string, 0, len(customerIDSet))
	for id := range customerIDSet {
		customerIDs = append(customerIDs, id)
	}

	customers, err := s.customerRepo.GetByIDs(ctx, customerIDs)
	if err != nil {
		return dto.ListPostsOutput{}, err
	}

	customerMap := make(map[string]entity.Customer, len(customers))
	for _, c := range customers {
		customerMap[c.ID] = c
	}

	for i := range out.Posts {
		post := &out.Posts[i]

		postMentions := mentionsMap[post.ID]

		post.Mentions = make([]entity.Customer, 0, len(postMentions))

		for _, m := range postMentions {
			if c, ok := customerMap[m.CustomerID]; ok {
				post.Mentions = append(post.Mentions, c)
			}
		}
	}

	return out, nil
}
