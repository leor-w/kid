package express

type Options struct {
	Key       string
	Customer  string
	Secret    string
	BaseUrl   string
	Salt      string
	NotifyUrl string
}

func WithKey(key string) Option {
	return func(options *Options) {
		options.Key = key
	}
}

func WithCustomer(customer string) Option {
	return func(options *Options) {
		options.Customer = customer
	}
}

func WithSecret(secret string) Option {
	return func(options *Options) {
		options.Secret = secret
	}
}

func WithBaseUrl(baseUrl string) Option {
	return func(options *Options) {
		options.BaseUrl = baseUrl
	}
}

func WithSalt(salt string) Option {
	return func(options *Options) {
		options.Salt = salt
	}
}

func WithNotifyUrl(notifyUrl string) Option {
	return func(options *Options) {
		options.NotifyUrl = notifyUrl
	}
}
