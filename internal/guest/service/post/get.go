package post

import (
	"context"
	"fmt"

	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
)

func (s *service) Get(ctx context.Context, postID string, reqCustomerID string) (entity.Post, error) {
	post, err := s.postRepo.Get(ctx, postID, reqCustomerID)
	if err != nil {
		return entity.Post{}, fmt.Errorf("get post from post repository: %w", err)
	}

	mentionsMap, err := s.postRepo.ListMentionsByPostIDs(ctx, []string{post.ID})
	if err != nil {
		return entity.Post{}, err
	}

	postMentions := mentionsMap[post.ID]

	if len(postMentions) == 0 {
		return post, nil
	}

	customerIDs := make([]string, 0, len(postMentions))
	for _, m := range postMentions {
		customerIDs = append(customerIDs, m.CustomerID)
	}

	customers, err := s.customerRepo.GetByIDs(ctx, customerIDs)
	if err != nil {
		return entity.Post{}, err
	}

	customerMap := make(map[string]entity.Customer, len(customers))
	for _, c := range customers {
		customerMap[c.ID] = c
	}

	post.Mentions = make([]entity.Customer, 0, len(postMentions))

	for _, m := range postMentions {
		if c, ok := customerMap[m.CustomerID]; ok {
			post.Mentions = append(post.Mentions, c)
		}
	}

	return post, nil
}
