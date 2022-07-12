package lock

import "time"

type Lock interface {
	Init(...Option)
	Lock(key string, ttl ...time.Duration) error
	Check(key string) (bool, error)
	Unlock(key string) error
	TTL(key string) (int64, error)
}

type Option func(*Options)
