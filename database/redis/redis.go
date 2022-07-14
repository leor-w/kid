package redis

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/leor-w/kid/config"
	"time"
)

type Client struct {
	*redis.Client
	options *Options
}

func (cli *Client) Provide() interface{} {
	return New(
		WithHost(config.GetString("redis.host")),
		WithPort(config.GetInt("redis.port")),
		WithDb(config.GetInt("redis.db")),
		WithPassword(config.GetString("redis.password")),
		WithDialTimeout(config.GetDuration("redis.dialTimeout")),
		WithReadTimeout(config.GetDuration("redis.readTimeout")),
		WithWriteTimeout(config.GetDuration("redis.writeTimeout")),
		WithMaxConnAge(config.GetDuration("redis.maxConnAge")),
		WithPoolTimeout(config.GetDuration("redis.poolTimeout")),
		WithIdleTimeout(config.GetDuration("redis.idleTimeout")),
		WithPoolSize(config.GetInt("redis.poolSize")),
		WithMinIdle(config.GetInt("redis.minIdleConn")),
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
		Addr:         fmt.Sprintf("%s:%d", cli.options.Host, cli.options.Port),
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
