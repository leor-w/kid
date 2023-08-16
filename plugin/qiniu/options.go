package qiniu

import (
	"time"
)

type Options struct {
	domain        string
	bucket        string
	access        string
	secret        string
	bucketConf    *Config
	operationConf *Config
	private       bool
	ttl           time.Duration
}

type Config struct {
	UseHttps      bool
	UseCdnDomains bool
	CentralRsHost string
	NotifyUrl     string
	Pipeline      string
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

func WithStorageConfig(config *Config) Option {
	return func(o *Options) {
		o.bucketConf = config
	}
}

func WithOperationConf(conf *Config) Option {
	return func(o *Options) {
		o.operationConf = conf
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
