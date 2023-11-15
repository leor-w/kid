package redis

import "time"

type Repository struct {
	RDB *Client `inject:""`
}

func (repo *Repository) Set(key string, value interface{}, expire int64) error {
	return repo.RDB.Set(key, value, time.Duration(expire)).Err()
}

func (repo *Repository) Get(key string) (string, error) {
	return repo.RDB.Get(key).Result()
}

func (repo *Repository) Del(keys ...string) error {
	return repo.RDB.Del(keys...).Err()
}

func (repo *Repository) Expire(key string, expire int64) error {
	return repo.RDB.Expire(key, time.Duration(expire)).Err()
}

func (repo *Repository) Exists(key string) (bool, error) {
	exist, err := repo.RDB.Exists(key).Result()
	if err != nil {
		return false, err
	}
	if exist > 0 {
		return true, nil
	}
	return false, nil
}
