package redis

import (
	"context"
	"fmt"

	"github.com/lsjhtang/kid/injector"

	"github.com/go-redis/redis/v8"
	"github.com/lsjhtang/kid/config"
)

type ClusterClient struct {
	*redis.ClusterClient
}

func (cli *ClusterClient) Provide(ctx context.Context) interface{} {
	var confName string
	name, ok := ctx.Value(injector.NameKey{}).(string)
	if ok && len(name) > 0 {
		confName = "." + name
	}
	confPrefix := fmt.Sprintf("redisCluster%s", confName)
	if !config.Exist(confPrefix) {
		panic(fmt.Sprintf("配置文件中缺少配置项: [%s]", confPrefix))
	}
	return NewCluster()
}

func NewCluster() interface{} {
	//rdb := redis.NewClusterClient(&redis.ClusterOptions{
	//	Addrs: []string{},
	//	NewClient: func(opt *redis.Options) *redis.Client {
	//
	//	},
	//	MaxRedirects:       0,
	//	ReadOnly:           false,
	//	RouteByLatency:     false,
	//	RouteRandomly:      false,
	//	ClusterSlots:       nil,
	//	Dialer:             nil,
	//	OnConnect:          nil,
	//	Username:           "",
	//	Password:           "",
	//	MaxRetries:         0,
	//	MinRetryBackoff:    0,
	//	MaxRetryBackoff:    0,
	//	DialTimeout:        0,
	//	ReadTimeout:        0,
	//	WriteTimeout:       0,
	//	PoolFIFO:           false,
	//	PoolSize:           0,
	//	MinIdleConns:       0,
	//	MaxConnAge:         0,
	//	PoolTimeout:        0,
	//	IdleTimeout:        0,
	//	IdleCheckFrequency: 0,
	//	TLSConfig:          nil,
	//},
	//)
	return nil
}
