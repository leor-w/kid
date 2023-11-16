package redis

import (
	"context"
	"time"
)

type RedisRepository struct {
	RDB *Client `inject:""`
}

func (repo *RedisRepository) Provide(context.Context) any {
	return repo
}

func (repo *RedisRepository) Set(key string, value interface{}, expire int64) error {
	return repo.RDB.Set(key, value, time.Duration(expire)*time.Second).Err()
}

func (repo *RedisRepository) Get(key string) (string, error) {
	return repo.RDB.Get(key).Result()
}

func (repo *RedisRepository) Del(keys ...string) error {
	return repo.RDB.Del(keys...).Err()
}

func (repo *RedisRepository) Expire(key string, expire int64) error {
	return repo.RDB.Expire(key, time.Duration(expire)*time.Second).Err()
}

func (repo *RedisRepository) Exists(key string) (bool, error) {
	exist, err := repo.RDB.Exists(key).Result()
	if err != nil {
		return false, err
	}
	if exist > 0 {
		return true, nil
	}
	return false, nil
}
