package token

import "time"

type Options struct {
	secret []byte
	expire time.Duration
}

func WithSecret(secret string) Option {
	return func(o *Options) {
		o.secret = []byte(secret)
	}
}

func WithExpire(expire int) Option {
	return func(o *Options) {
		o.expire = time.Duration(expire) * time.Hour * 24
	}
}
