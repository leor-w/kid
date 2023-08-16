package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/leor-w/kid/config"
	"github.com/leor-w/kid/plugin"
	"github.com/leor-w/kid/utils"
)

type Client struct {
	*redis.Client
	options *Options
}

func (cli *Client) Provide(ctx context.Context) interface{} {
	var confName string
	name, ok := ctx.Value(plugin.NameKey{}).(string)
	if ok && len(name) > 0 {
		confName = "." + name
	}
	confPrefix := fmt.Sprintf("redis%s", confName)
	if !config.Exist(confPrefix) {
		panic(fmt.Sprintf("config file not found configuration item [%s]", confPrefix))
	}
	return New(
		WithAddr(config.GetString(utils.GetConfigurationItem(confPrefix, "addr"))),
		WithDb(config.GetInt(utils.GetConfigurationItem(confPrefix, "db"))),
		WithPassword(config.GetString(utils.GetConfigurationItem(confPrefix, "password"))),
		WithDialTimeout(config.GetDuration(utils.GetConfigurationItem(confPrefix, "dialTimeout"))),
		WithReadTimeout(config.GetDuration(utils.GetConfigurationItem(confPrefix, "readTimeout"))),
		WithWriteTimeout(config.GetDuration(utils.GetConfigurationItem(confPrefix, "writeTimeout"))),
		WithMaxConnAge(config.GetDuration(utils.GetConfigurationItem(confPrefix, "maxConnAge"))),
		WithPoolTimeout(config.GetDuration(utils.GetConfigurationItem(confPrefix, "poolTimeout"))),
		WithIdleTimeout(config.GetDuration(utils.GetConfigurationItem(confPrefix, "idleTimeout"))),
		WithCtxTimeout(config.GetDuration(utils.GetConfigurationItem(confPrefix, "ctxTimeout"))),
		WithPoolSize(config.GetInt(utils.GetConfigurationItem(confPrefix, "poolSize"))),
		WithMinIdle(config.GetInt(utils.GetConfigurationItem(confPrefix, "minIdleConn"))),
	)
}

type Option func(*Options)

func New(opts ...Option) *Client {
	var options = &Options{
		DbNum:        0,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		MaxConnAge:   2 * time.Hour,
		PoolTimeout:  5 * time.Second,
		IdleTimeout:  10 * time.Minute,
		PoolSize:     50,
		MinIdleConn:  10,
	}

	for _, o := range opts {
		o(options)
	}
	cli := &Client{
		options: options,
	}
	cli.Client = redis.NewClient(&redis.Options{
		Addr:         cli.options.Addr,
		Password:     cli.options.Password,
		DB:           cli.options.DbNum,
		DialTimeout:  cli.options.DialTimeout,
		ReadTimeout:  cli.options.ReadTimeout,
		WriteTimeout: cli.options.WriteTimeout,
		PoolSize:     cli.options.PoolSize,
		MinIdleConns: cli.options.MinIdleConn,
		MaxConnAge:   cli.options.MaxConnAge,
		PoolTimeout:  cli.options.PoolTimeout,
		IdleTimeout:  cli.options.IdleTimeout,
	})
	if err := cli.Ping().Err(); err != nil {
		panic(fmt.Sprintf("ping redis failed: %s", err.Error()))
	}
	return cli
}
