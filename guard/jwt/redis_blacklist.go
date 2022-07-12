package jwt

import (
	"fmt"
	"github.com/leor-w/kid/database/redis"
	"time"
)

type RedisBlacklist struct {
	conn redis.Conn
}

func (rbl *RedisBlacklist) Black(license string, ttl time.Duration) error {
	return rbl.conn.Set(fmt.Sprintf(BlacklistKey, license), "", ttl).Err()
}

func (rbl *RedisBlacklist) Checklist(license string) (bool, error) {
	exists, err := rbl.conn.Exists(fmt.Sprintf(BlacklistKey, license)).Result()
	if err != nil {
		return true, err
	}
	if exists > 0 {
		return true, nil
	}
	return false, nil
}
