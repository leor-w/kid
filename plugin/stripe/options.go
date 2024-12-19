package stripe

type Options struct {
	SecretKey      string
	WebhookSecret  string
	RedirectType   string
	RedirectDomain string
}

func WithSecretKey(secretKey string) Option {
	return func(o *Options) {
		o.SecretKey = secretKey
	}
}

func WithWebhookSecret(webhookSecret string) Option {
	return func(o *Options) {
		o.WebhookSecret = webhookSecret
	}
}

func WithRedirectType(redirectType string) Option {
	return func(options *Options) {
		options.RedirectType = redirectType
	}
}

func WithRedirectDomain(redirectDomain string) Option {
	return func(o *Options) {
		o.RedirectDomain = redirectDomain
	}
}
