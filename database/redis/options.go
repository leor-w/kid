package redis

import "time"

type Options struct {
	Addr         string
	DbNum        int
	Password     string
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	MaxConnAge   time.Duration
	PoolTimeout  time.Duration
	IdleTimeout  time.Duration
	CtxTimeout   time.Duration
	PoolSize     int
	MinIdleConn  int
}

func WithAddr(addr string) Option {
	return func(o *Options) {
		o.Addr = addr
	}
}

func WithPassword(password string) Option {
	return func(o *Options) {
		o.Password = password
	}
}

func WithDb(dbNum int) Option {
	return func(o *Options) {
		o.DbNum = dbNum
	}
}

func WithDialTimeout(dial time.Duration) Option {
	return func(o *Options) {
		o.DialTimeout = dial * time.Second
	}
}

func WithReadTimeout(read time.Duration) Option {
	return func(o *Options) {
		o.ReadTimeout = read * time.Second
	}
}

func WithWriteTimeout(write time.Duration) Option {
	return func(o *Options) {
		o.WriteTimeout = write * time.Second
	}
}

func WithPoolTimeout(pool time.Duration) Option {
	return func(o *Options) {
		o.PoolTimeout = pool * time.Second
	}
}

func WithIdleTimeout(idle time.Duration) Option {
	return func(o *Options) {
		o.IdleTimeout = idle * time.Minute
	}
}

func WithCtxTimeout(ctxTimeout time.Duration) Option {
	return func(o *Options) {
		o.CtxTimeout = ctxTimeout * time.Second
	}
}

func WithPoolSize(poolSize int) Option {
	return func(o *Options) {
		o.PoolSize = poolSize
	}
}

func WithMinIdle(minIdle int) Option {
	return func(o *Options) {
		o.MinIdleConn = minIdle
	}
}

func WithMaxConnAge(maxConnAge time.Duration) Option {
	return func(o *Options) {
		o.MaxConnAge = maxConnAge * time.Hour
	}
}
