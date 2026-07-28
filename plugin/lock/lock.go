package lock

import (
	"context"
	"time"

	"github.com/lsjhtang/kid/logger"

	"github.com/lsjhtang/kid/database/redis"
)

type Lock interface {
	Lock(key string, ttl time.Duration) bool
	Unlock(key string) error
}

type Option func(*Options)

type RedisLock struct {
	rdb redis.Conn `inject:""`
}

func (rl *RedisLock) Provide(context.Context) interface{} {
	return new(RedisLock)
}

func (rl *RedisLock) Lock(key string, ttl time.Duration) bool {
	result, err := rl.rdb.SetNX(key, 1, ttl).Result()
	if err != nil {
		logger.Errorf("lock: redis 锁获取失败: %s", err.Error())
		return false
	}
	return result
}

func (rl *RedisLock) Check(key string) (bool, error) {
	exist, err := rl.rdb.Exists(key).Result()
	if err != nil {
		return false, err
	}
	if exist > 0 {
		return true, nil
	}
	return false, nil
}

func (rl *RedisLock) Unlock(key string) error {
	return rl.rdb.Expire(key, 0).Err()
}
