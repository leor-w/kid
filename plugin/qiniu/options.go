package qiniu

import (
	"github.com/qiniu/go-sdk/v7/storage"
	"time"
)

type Options struct {
	domain     string
	bucket     string
	access     string
	secret     string
	bucketConf *storage.Config
	private    bool
	ttl        time.Duration
}

func WithDomain(domain string) Option {
	return func(o *Options) {
		o.domain = domain
	}
}

func WithBucket(bucket string) Option {
	return func(o *Options) {
		o.bucket = bucket
	}
}

func WithAccess(access string) Option {
	return func(o *Options) {
		o.access = access
	}
}

func WithSecret(secret string) Option {
	return func(o *Options) {
		o.secret = secret
	}
}

func WithStorageConfig(config *storage.Config) Option {
	return func(o *Options) {
		o.bucketConf = config
	}
}

func WithPrivate(private bool) Option {
	return func(o *Options) {
		o.private = private
	}
}

func WithTTL(ttl time.Duration) Option {
	return func(o *Options) {
		o.ttl = ttl
	}
}
