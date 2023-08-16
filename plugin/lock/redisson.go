package lock

import "github.com/go-redis/redis/v8"

type IRedisson interface {
}

type Redisson struct {
	rds *redis.Client        `inject:"NR"`
	rcc *redis.ClusterClient `inject:"NR"`
}
