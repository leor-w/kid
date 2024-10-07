package tiktok

type Options struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scope        []string
}

func WithClientID(clientID string) Option {
	return func(o *Options) {
		o.ClientID = clientID
	}
}

func WithClientSecret(clientSecret string) Option {
	return func(o *Options) {
		o.ClientSecret = clientSecret
	}
}

func WithRedirectURL(redirectURL string) Option {
	return func(o *Options) {
		o.RedirectURL = redirectURL
	}
}

func WithScope(scope []string) Option {
	return func(o *Options) {
		o.Scope = scope
	}
}
