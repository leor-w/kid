package lock

import (
	"github.com/leor-w/kid/database/redis"
	"time"
)

type Lock interface {
	Init(...Option)
	Lock(key string, ttl ...time.Duration) error
	Check(key string) (bool, error)
	Unlock(key string) error
	TTL(key string) (int64, error)
}

type Option func(*Options)

type RedisLock struct {
	options *Options
	rdb     redis.Conn
}

func (rl *RedisLock) Provide() interface{} {
	return nil
}

func (rl *RedisLock) Init(opts ...Option) {
}

func (rl *RedisLock) Lock(key string, ttl ...time.Duration) error {
	return nil
}

func (rl *RedisLock) Check(key string) (bool, error) {
	return false, nil
}

func (rl *RedisLock) Unlock(key string) error {
	return nil
}

func (rl *RedisLock) TTL(key string) (int64, error) {
	return 0, nil
}
