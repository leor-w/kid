package stripe

type Options struct {
	SecretKey     string
	WebhookSecret string
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
