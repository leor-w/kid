package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	redis2 "github.com/go-redis/redis/v8"

	"github.com/lsjhtang/kid/database/redis"
	"github.com/lsjhtang/kid/guard"
	"github.com/lsjhtang/kid/guard/constant"
)

type RedisStore struct {
	rdb redis.Conn `inject:""`
}

func (rs *RedisStore) Provide(ctx context.Context) interface{} {
	return &RedisStore{}
}

// GetTokenInfo 通过 token 获取 token 详情
func (rs *RedisStore) GetTokenInfo(license string) (*guard.TokenInfo, error) {
	val, err := rs.rdb.Get(constant.GetTokenDetailKey(license)).Result()
	if err != nil {
		if errors.Is(err, redis2.Nil) {
			return nil, fmt.Errorf("token [%s] not found", license)
		}
		return nil, err
	}
	var tokenDetail guard.TokenInfo
	if err := json.Unmarshal([]byte(val), &tokenDetail); err != nil {
		return nil, err
	}
	return &tokenDetail, nil
}

// GetUserTokens 通过用户类型及用户的 uid 获取用户的 token 列表
func (rs *RedisStore) GetUserTokens(uType guard.UserType, uid int64) ([]string, error) {
	var (
		offset uint64
		tokens []string
	)
	// 通过用户类型及用户的 uid 获取用户的 token 列表
	for {
		keys, cursor, err := rs.rdb.Scan(context.Background(), offset, constant.GetUserTokenSearchKey(uType, uid), 10).Result()
		if err != nil {
			return nil, err
		}
		for _, key := range keys {
			tokens = append(tokens, constant.GetTokenByKey(key))
		}
		if cursor == 0 {
			break
		}
		offset = cursor
	}
	return tokens, nil
}

// Save 保存用户的 token
func (rs *RedisStore) Save(license string, tokenInfo *guard.TokenInfo) error {
	infoBytes, err := json.Marshal(tokenInfo)
	if err != nil {
		return err
	}
	expired := time.Duration(tokenInfo.ExpiredAt-time.Now().Unix()) * time.Second
	if _, err := rs.rdb.Pipelined(func(pipe redis2.Pipeliner) error {
		// 保存用户token
		if err := pipe.Set(context.Background(), constant.GetUserTokensKey(tokenInfo.Type, tokenInfo.Uid, license), string(infoBytes), expired).Err(); err != nil {
			return err
		}
		// 保存token详情
		if err := pipe.Set(context.Background(), constant.GetTokenDetailKey(license), string(infoBytes), expired).Err(); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

// Expired 使 token 失效
func (rs *RedisStore) Expired(license string) error {
	tokenInfo, err := rs.GetTokenInfo(license)
	if err != nil {
		return err
	}
	if _, err := rs.rdb.Pipelined(func(pipe redis2.Pipeliner) error {
		// 删除用户token
		if err := pipe.Expire(context.Background(), constant.GetUserTokensKey(tokenInfo.Type, tokenInfo.Uid, license), 0).Err(); err != nil {
			return err
		}
		// 删除token详情
		if err = pipe.Expire(context.Background(), constant.GetTokenDetailKey(license), 0).Err(); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (rs *RedisStore) ExpiredAll(uType guard.UserType, uid int64) error {
	userTokens, err := rs.GetUserTokens(uType, uid)
	if err != nil {
		return err
	}
	for _, token := range userTokens {
		if err := rs.Expired(token); err != nil {
			return err
		}
	}
	return nil
}

// Exist 判断 token 是否存在
func (rs *RedisStore) Exist(license string) bool {
	exist, err := rs.rdb.Exists(constant.GetTokenDetailKey(license)).Result()
	if err != nil {
		return false
	}
	if exist <= 0 {
		return false
	}
	tokenInfo, err := rs.GetTokenInfo(license)
	if err != nil {
		return false
	}
	if tokenInfo.ExpiredAt < time.Now().Unix() {
		return false
	}
	if rs.rdb.Exists(constant.GetUserTokensKey(tokenInfo.Type, tokenInfo.Uid, license)).Val() <= 0 {
		return false
	}
	return true
}
