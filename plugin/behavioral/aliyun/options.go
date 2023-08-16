package aliyun

type Options struct {
	AccessKeyId     string
	AccessKeySecret string
	RegionId        string
}

type Option func(*Options)

func WithAccessKeyId(keyId string) Option {
	return func(o *Options) {
		o.AccessKeyId = keyId
	}
}

func WithAccessKeySecret(secret string) Option {
	return func(o *Options) {
		o.AccessKeySecret = secret
	}
}

func WithRegionId(regionId string) Option {
	return func(o *Options) {
		o.RegionId = regionId
	}
}
