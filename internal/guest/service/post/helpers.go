package post

import (
	"context"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	"github.com/ua-academy-projects/share-bite/internal/storage"
	"github.com/ua-academy-projects/share-bite/pkg/logger"
)

func rollbackUploadedImages(objectStorage storage.ObjectStorage, keys []string) {
	if objectStorage == nil || len(keys) == 0 {
		return
	}

	for _, key := range keys {
		cleanupDelete(
			objectStorage,
			key,
		)
	}
}

func cleanupDelete(objectStorage storage.ObjectStorage, key string) {
	if objectStorage == nil || key == "" {
		return
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		cleanupTimeout,
	)
	defer cancel()

	if err := objectStorage.Delete(
		ctx,
		key,
	); err != nil {

		logger.WarnKV(
			ctx,
			"failed to cleanup post image object",
			"key",
			key,
			"error",
			err,
		)
	}
}

func UniqueStrings(input []string) []string {
	set := make(map[string]struct{}, len(input))
	res := make([]string, 0, len(input))

	for _, value := range input {
		if _, ok := set[value]; ok {
			continue
		}

		set[value] = struct{}{}
		res = append(res, value)
	}

	return res
}

func uniqueAndExcludeSelf(self string, ids []string) []string {
	m := make(map[string]struct{})

	for _, id := range ids {
		if id == self {
			continue
		}
		m[id] = struct{}{}
	}

	res := make([]string, 0, len(m))
	for id := range m {
		res = append(res, id)
	}

	return res
}

func isValidPostStatusTransition(from, to entity.PostStatus) bool {
	switch from {
	case entity.PostStatusDraft:
		return to == entity.PostStatusDraft ||
			to == entity.PostStatusPublished ||
			to == entity.PostStatusArchived
	case entity.PostStatusPublished:
		return to == entity.PostStatusPublished ||
			to == entity.PostStatusArchived
	case entity.PostStatusArchived:
		return to == entity.PostStatusArchived ||
			to == entity.PostStatusPublished
	default:
		return false
	}
}
