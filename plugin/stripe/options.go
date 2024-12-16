package stripe

type Options struct {
	SecretKey     string
	WebhookSecret string
	RedirectType  string
	RedirectURL   string
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

func WithRedirectURL(redirectURL string) Option {
	return func(o *Options) {
		o.RedirectURL = redirectURL
	}
}
