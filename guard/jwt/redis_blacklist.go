package jwt

import (
	"context"
	"time"

	"github.com/lsjhtang/kid/database/redis"
)

type RedisBlacklist struct {
	conn redis.Conn `inject:""`
}

func (rbl *RedisBlacklist) Provide(ctx context.Context) interface{} {
	return &RedisBlacklist{}
}

func (rbl *RedisBlacklist) Black(license string, ttl time.Duration) error {
	return rbl.conn.Set(GetBlacklistKey(license), license, ttl).Err()
}

func (rbl *RedisBlacklist) IsBlacked(license string) bool {
	exists, err := rbl.conn.Exists(GetBlacklistKey(license)).Result()
	if err != nil {
		return true
	}
	if exists > 0 {
		return true
	}
	return false
}
