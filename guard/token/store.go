package token

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/leor-w/kid/config"
	"github.com/leor-w/kid/database/redis"
	"github.com/leor-w/kid/guard"
	"time"
)

type Store interface {
	Get(license string) (*tokenInfo, error)
	Save(license string, user *tokenInfo) error
	Expired(license string) error
	Exist(token string) bool
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
	var info tokenInfo
	err = json.Unmarshal([]byte(val), &info)
	if err != nil {
		return nil, err
	}
	exist, err = rs.rdb.Exists(GetTokenIdKey(info.Type, info.Id)).Result()
	if err != nil {
		return nil, err
	}
	if exist <= 0 {
		return nil, errors.New("token bind user not exist")
	}
	redisToken, err := rs.rdb.Get(GetTokenIdKey(info.Type, info.Id)).Result()
	if err != nil {
		return nil, err
	}
	if redisToken != token {
		return nil, errors.New("request token invalid")
	}
	return &info, nil
}

func (rs *RedisStore) Save(token string, session *tokenInfo) error {
	if err := rs.removeOldToken(session.Type, session.Id); err != nil {
		return err
	}
	valBytes, err := json.Marshal(session)
	if err != nil {
		return err
	}
	expired := time.Duration(session.ExpireAt-time.Now().Unix()) * time.Second
	if err = rs.rdb.Set(GetTokenIdKey(session.Type, session.Id), session.Token, expired).Err(); err != nil {
		return err
	}
	return rs.rdb.Set(GetTokenKey(token), string(valBytes), expired).Err()
}

func (rs *RedisStore) Expired(license string) error {
	info, err := rs.Get(license)
	if err != nil {
		return err
	}
	if err = rs.rdb.Expire(GetTokenIdKey(info.Type, info.Id), 0).Err(); err != nil {
		return err
	}
	return rs.rdb.Expire(GetTokenKey(license), 0).Err()
}

func (rs *RedisStore) removeOldToken(userType guard.UserType, tokenUid int64) error {
	exist, err := rs.rdb.Exists(GetTokenIdKey(userType, tokenUid)).Result()
	if err != nil {
		return err
	}
	if exist <= 0 {
		return nil
	}
	oldToken, err := rs.rdb.Get(GetTokenIdKey(userType, tokenUid)).Result()
	if err != nil {
		return err
	}
	if err = rs.rdb.Expire(GetTokenKey(oldToken), 0).Err(); err != nil {
		return err
	}
	if err = rs.rdb.Expire(GetTokenIdKey(userType, tokenUid), 0).Err(); err != nil {
		return err
	}
	return nil
}

func (rs *RedisStore) Exist(token string) bool {
	exist, err := rs.rdb.Exists(GetTokenKey(token)).Result()
	if err != nil {
		return true
	}
	if exist > 0 {
		return true
	}
	return false
}
