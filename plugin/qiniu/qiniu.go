package qiniu

import (
	"github.com/leor-w/kid/config"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

type Qiniu struct {
	options       *Options
	mac           *qbox.Mac
	bucketManager *storage.BucketManager
}

func (qiniu *Qiniu) Provide() interface{} {
	if !config.Exist("qiniu") {
		panic("not found [qiniu] in config")
	}
	return New(
		WithDomain(config.GetString("qiniu.domain")),
		WithBucket(config.GetString("qiniu.bucket")),
		WithAccess(config.GetString("qiniu.access")),
		WithSecret(config.GetString("qiniu.secret")),
		WithPrivate(config.GetBool("qiniu.private")),
		WithTTL(config.GetDuration("qiniu.ttl")),
	)
}

type Option func(*Options)

func (qiniu *Qiniu) UploadToken(keys ...string) string {
	if len(keys) <= 0 {
		policy := storage.PutPolicy{Scope: qiniu.options.bucket}
		return policy.UploadToken(qiniu.mac)
	}
	policy := storage.PutPolicy{Scope: qiniu.options.bucket + ":" + keys[0]}
	return policy.UploadToken(qiniu.mac)
}

func New(opts ...Option) *Qiniu {
	options := &Options{
		private: false,
		ttl:     5,
	}
	for _, opt := range opts {
		opt(options)
	}
	mac := qbox.NewMac(options.access, options.secret)
	return &Qiniu{
		options:       options,
		mac:           mac,
		bucketManager: storage.NewBucketManager(mac, options.bucketConf),
	}
}
