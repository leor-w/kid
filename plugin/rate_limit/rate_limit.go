package rate_limit

import (
	"context"
	"fmt"

	redis2 "github.com/go-redis/redis/v8"
	"github.com/lsjhtang/kid/injector"

	"github.com/lsjhtang/kid/config"
	"github.com/lsjhtang/kid/database/redis"
	"github.com/lsjhtang/kid/utils"
)

type RateLimit struct {
	rdb     *redis.Client `inject:""`
	options *Options
}

type Option func(*Options)

func (r *RateLimit) Provide(ctx context.Context) any {
	var confName string
	if name, ok := ctx.Value(injector.NameKey{}).(string); ok && len(name) > 0 {
		confName = "." + name
	}
	confPrefix := fmt.Sprintf("rate_limit%s", confName)
	if !config.Exist(confPrefix) {
		panic(fmt.Sprintf("配置文件未找到配置项 [%s]", confPrefix))
	}
	return New(
		WithRateLimit(config.GetInt(utils.GetConfigurationItem(confPrefix, "rate_limit"))),
		WithBurstLimit(config.GetInt(utils.GetConfigurationItem(confPrefix, "burst_limit"))),
		WithRateLimitWindow(config.GetInt(utils.GetConfigurationItem(confPrefix, "rate_limit_window"))),
	)
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

func New(options ...Option) *RateLimit {
	var opts = &Options{
		RateLimit:       5,
		BurstLimit:      10,
		RateLimitWindow: 60,
	}
	for _, option := range options {
		option(opts)
	}
	return &RateLimit{
		options: opts,
	}
}
