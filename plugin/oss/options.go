package oss

type Options struct {
	endpoint   string
	bucketName string
	accessKey  string
	secretKey  string
}

func WithEndpoint(endpoint string) Option {
	return func(o *Options) {
		o.endpoint = endpoint
	}
}

func WithBucketName(bucketName string) Option {
	return func(o *Options) {
		o.bucketName = bucketName
	}
}

func WithAccessKey(accessKey string) Option {
	return func(o *Options) {
		o.accessKey = accessKey
	}
}

func WithSecretKey(secretKey string) Option {
	return func(o *Options) {
		o.secretKey = secretKey
	}
}
