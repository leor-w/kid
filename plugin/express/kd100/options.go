package kd100

type Options struct {
	Key                  string
	Customer             string
	Secret               string
	BaseUrl              string
	Salt                 string
	MapTempKey           string
	NotifyUrl            string
	SendExpressNotifyUrl string
}

func WithKey(key string) Option {
	return func(o *Options) {
		o.Key = key
	}
}

func WithCustomer(customer string) Option {
	return func(o *Options) {
		o.Customer = customer
	}
}

func WithSecret(secret string) Option {
	return func(o *Options) {
		o.Secret = secret
	}
}

func WithBaseUrl(baseUrl string) Option {
	return func(o *Options) {
		o.BaseUrl = baseUrl
	}
}

func WithSalt(salt string) Option {
	return func(o *Options) {
		o.Salt = salt
	}
}

func WithMapTempKey(mapTempKey string) Option {
	return func(o *Options) {
		o.MapTempKey = mapTempKey
	}
}

func WithNotifyUrl(notifyUrl string) Option {
	return func(o *Options) {
		o.NotifyUrl = notifyUrl
	}
}

func WithSendExpressNotifyUrl(notifyUrl string) Option {
	return func(o *Options) {
		o.SendExpressNotifyUrl = notifyUrl
	}
}
