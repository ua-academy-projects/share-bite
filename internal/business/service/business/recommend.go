package business

import (
	"context"
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

func (s *service) RecommendVenues(ctx context.Context, userID string, skip, limit int) (pagination.Result[entity.OrgUnit], error) {
	const op = "service.business.RecommendVenues"

	const tagsToFetch = 5

	topTags, err := s.businessRepo.GetTopTagsByUserLikes(ctx, userID, tagsToFetch)
	if err != nil {
		return pagination.Result[entity.OrgUnit]{}, err
	}

	quotas := calculateTagQuotas(topTags, limit)

	var finalVenues []entity.OrgUnit
	var seenIDs []int
	deficit := 0

	for _, tag := range topTags {
		quota := quotas[tag]
		if quota == 0 {
			continue
		}

		venues, err := s.businessRepo.GetVenuesByTag(ctx, tag, quota, seenIDs)
		if err != nil {
			return pagination.Result[entity.OrgUnit]{}, err
		}

		finalVenues = append(finalVenues, venues...)
		for _, v := range venues {
			seenIDs = append(seenIDs, v.Id)
		}

		if len(venues) < quota {
			deficit += quota - len(venues)
		}
	}

	if deficit > 0 {
		fillVenues, err := s.businessRepo.GetRandomVenues(ctx, deficit, seenIDs)
		if err == nil {
			finalVenues = append(finalVenues, fillVenues...)
		}
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(finalVenues), func(i, j int) {
		finalVenues[i], finalVenues[j] = finalVenues[j], finalVenues[i]
	})

	return pagination.Result[entity.OrgUnit]{
		Items: finalVenues,
		Total: len(finalVenues),
	}, nil
}
