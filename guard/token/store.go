package token

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/leor-w/kid/config"
	"github.com/leor-w/kid/database/redis"
	"time"
)

type Store interface {
	Get(license string) (*tokenInfo, error)
	Save(license string, user *tokenInfo) error
	Expired(license string) error
}

type RedisStore struct {
	rdb redis.Conn `inject:""`
}

func (rs *RedisStore) Provide() interface{} {
	return NewStore(config.DefaultString("token.store", "redis"))
}

func NewStore(cacheType string) Store {
	switch cacheType {
	case "redis":
		return &RedisStore{}
	default:
		return nil
	}
}

func (rs *RedisStore) Get(token string) (*tokenInfo, error) {
	exist, err := rs.rdb.Exists(GetTokenKey(token)).Result()
	if err != nil {
		return nil, err
	}
	if exist <= 0 {
		return nil, fmt.Errorf("token not exist")
	}
	val, err := rs.rdb.Get(GetTokenKey(token)).Result()
	if err != nil {
		return nil, err
	}
	var sessionInfo tokenInfo
	err = json.Unmarshal([]byte(val), &sessionInfo)
	if err != nil {
		return nil, err
	}
	exist, err = rs.rdb.Exists(GetTokenIdKey(sessionInfo.Id)).Result()
	if err != nil {
		return nil, err
	}
	if exist <= 0 {
		return nil, errors.New("token bind user not exist")
	}
	redisToken, err := rs.rdb.Get(GetTokenIdKey(sessionInfo.Id)).Result()
	if err != nil {
		return nil, err
	}
	if redisToken != token {
		return nil, errors.New("request token invalid")
	}
	return &sessionInfo, nil
}

func (rs *RedisStore) Save(token string, session *tokenInfo) error {
	valBytes, err := json.Marshal(session)
	if err != nil {
		return err
	}
	expired := time.Duration(session.ExpireAt-time.Now().Unix()) * time.Second
	if err = rs.rdb.Set(GetTokenIdKey(session.Id), session.Token, expired).Err(); err != nil {
		return err
	}
	return rs.rdb.Set(GetTokenKey(token), string(valBytes), expired).Err()
}

func (rs *RedisStore) Expired(license string) error {
	return rs.rdb.Expire(GetTokenKey(license), 0).Err()
}
