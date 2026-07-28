package captcha

import (
	"context"
	"time"

	"github.com/lsjhtang/kid/database/redis"
)

type RedisStore struct {
	rds *redis.Client `inject:""`
}

func (r *RedisStore) Provide(_ context.Context) interface{} {
	return r
}

func (r *RedisStore) Set(id string, value string) error {
	return r.rds.SetNX(id, value, time.Minute*3).Err()
}

func (r *RedisStore) Get(id string, clear bool) string {
	value, err := r.rds.Get(id).Result()
	if err != nil {
		return ""
	}
	if clear {
		r.rds.Expire(id, 0)
	}
	return value
}

func (r *RedisStore) Verify(id, answer string, clear bool) bool {
	value := r.Get(id, clear)
	if value == answer {
		return true
	}
	return false
}
