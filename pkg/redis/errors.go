package redis

import (
	"errors"

	goredis "github.com/redis/go-redis/v9"
)

func IsPermanentRedisError(err error) bool {
	if err == nil {
		return false
	}

	if errors.Is(err, goredis.ErrClosed) {
		return true
	}

	return goredis.IsAuthError(err) ||
		goredis.IsReadOnlyError(err) ||
		goredis.IsPermissionError(err)
}
