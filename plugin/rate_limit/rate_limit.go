package rate_limit

import (
	"context"
	"fmt"

	redis2 "github.com/go-redis/redis/v8"
	"github.com/leor-w/kid/database/redis"
)

type RateLimit struct {
	rdb     *redis.Client
	options *Options
}

func (r *RateLimit) RateLimit(key string) (bool, error) {
	redisKey := RateLimitIP + key

	// 使用 redis 的 Lua 脚本实现原子性操作限流
	luaScript := redis2.NewScript(LimitLuaScript)

	res, err := luaScript.Run(
		context.Background(),
		r.rdb.Client,
		[]string{redisKey},
		r.options.BurstLimit,
		r.options.RateLimitWindow,
	).Result()
	if err != nil {
		return false, fmt.Errorf("请求限流检查错误：%w", err)
	}
	if res == 0 {
		return false, nil
	}
	return true, nil
}

type Option func(*Options)
