package jwt

import (
	"time"
)

type Options struct {
	Issuer        string
	RefreshExpire time.Duration
	Expire        time.Duration
	SigningMethod SigningMethod
	Key           []byte
}

func WithIssuer(issuer string) Option {
	return func(o *Options) {
		o.Issuer = issuer
	}
}

func WithRefreshExpire(refreshExpire time.Duration) Option {
	return func(o *Options) {
		o.RefreshExpire = refreshExpire
	}
}

func WithExpire(expire time.Duration) Option {
	return func(o *Options) {
		o.Expire = expire
	}
}

func WithSigningMethod(signMethod SigningMethod) Option {
	return func(o *Options) {
		o.SigningMethod = signMethod
	}
}

func WithKey(key []byte) Option {
	return func(o *Options) {
		o.Key = key
	}
}
