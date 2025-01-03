package awss3

type Options struct {
	Region       string
	AccessKey    string
	SecretKey    string
	CacheControl string
}

func WithRegion(region string) Option {
	return func(o *Options) {
		o.Region = region
	}
}

func WithAccessKey(accessKey string) Option {
	return func(o *Options) {
		o.AccessKey = accessKey
	}
}

func WithSecretKey(secretKey string) Option {
	return func(o *Options) {
		o.SecretKey = secretKey
	}
}

func WithCacheControl(cacheControl string) Option {
	return func(o *Options) {
		o.CacheControl = cacheControl
	}
}
