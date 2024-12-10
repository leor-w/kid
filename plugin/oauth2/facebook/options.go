package facebook

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

func WithScope(scopes []string) Option {
	for _, s := range scopes {
		switch s {
		case "email":
			scopes = append(scopes, ScopeEmail)
		}
	}
	return func(o *Options) {
		o.Scope = scopes
	}
}
