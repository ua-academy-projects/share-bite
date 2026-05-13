package business

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/ua-academy-projects/share-bite/internal/business/entity"
	"github.com/ua-academy-projects/share-bite/pkg/database/pagination"
)

func calculateTagQuotas(tags []string, limit int) map[string]int {
	n := len(tags)
	if n == 0 || limit <= 0 {
		return nil
	}

	weights := make([]int, n)
	weights[n-1] = 1
	if n > 1 {
		weights[n-2] = 1
	}
	for i := n - 3; i >= 0; i-- {
		weights[i] = weights[i+1] + weights[i+2]
	}

	totalWeight := 0
	for _, w := range weights {
		totalWeight += w
	}

	quotas := make(map[string]int, n)
	allocated := 0

	for i, tag := range tags {
		q := (limit * weights[i]) / totalWeight
		quotas[tag] = q
		allocated += q
	}

	remainder := limit - allocated
	if remainder > 0 {
		for i := 0; i < remainder; i++ {
			quotas[tags[i%n]]++
		}
	}

	return quotas
}

func (s *service) RecommendPosts(ctx context.Context, userID string, lat, lon float64, skip, limit int) (pagination.Result[entity.RecomendedPost], error) {
	const op = "service.business.RecommendPosts"
	const tagsToFetch = 5

	h3Hashes := s.h3Service.GetH3Neighbors(lat, lon, s.h3Config.Resolution, s.h3Config.RecommendRadius)
	if len(h3Hashes) == 0 {
		return pagination.Result[entity.RecomendedPost]{}, nil
	}

	topTags, err := s.businessRepo.GetTopTagsByUserLikes(ctx, userID, tagsToFetch)
	if err != nil {
		return pagination.Result[entity.RecomendedPost]{}, err
	}

	if len(topTags) == 0 {
		fillPosts, err := s.businessRepo.GetRandomPosts(ctx, limit, nil, h3Hashes)
		if err != nil {
			return pagination.Result[entity.RecomendedPost]{}, err
		}

		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		r.Shuffle(len(fillPosts), func(i, j int) {
			fillPosts[i], fillPosts[j] = fillPosts[j], fillPosts[i]
		})

		return pagination.Result[entity.RecomendedPost]{
			Items: fillPosts,
			Total: len(fillPosts), // TODO: Get the count
		}, nil
	}

	quotas := calculateTagQuotas(topTags, limit)
	var finalPosts []entity.RecomendedPost
	var seenCompositeIDs []string
	deficit := 0

	for _, tag := range topTags {
		quota := quotas[tag]
		if quota == 0 {
			continue
		}

		posts, err := s.businessRepo.GetPostsByTag(ctx, tag, quota, seenCompositeIDs, h3Hashes)
		if err != nil {
			return pagination.Result[entity.RecomendedPost]{}, err
		}

		finalPosts = append(finalPosts, posts...)
		for _, p := range posts {
			seenCompositeIDs = append(seenCompositeIDs, fmt.Sprintf("%s:%d", p.PostType, p.ID))
		}

		if len(posts) < quota {
			deficit += quota - len(posts)
		}
	}

	if deficit > 0 {
		fillPosts, err := s.businessRepo.GetRandomPosts(ctx, deficit, seenCompositeIDs, h3Hashes)
		if err == nil {
			finalPosts = append(finalPosts, fillPosts...)
		}
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(finalPosts), func(i, j int) {
		finalPosts[i], finalPosts[j] = finalPosts[j], finalPosts[i]
	})

	return pagination.Result[entity.RecomendedPost]{
		Items: finalPosts,
		Total: len(finalPosts),
	}, nil
}
