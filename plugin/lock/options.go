package lock

import (
	"time"
)

type Options struct {
	lockTime time.Duration
}

func WithLockTime(lockTime time.Duration) Option {
	return func(o *Options) {
		o.lockTime = lockTime
	}
}
