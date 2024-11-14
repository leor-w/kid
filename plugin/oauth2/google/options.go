package google

import "time"

type Options struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scope        []string
	TTL          time.Duration
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
	scopes := make([]string, 0, len(scope))
	for _, s := range scope {
		switch s {
		case "email":
			scopes = append(scopes, Email)
		case "profile":
			scopes = append(scopes, Profile)
		}
	}
	return func(o *Options) {
		o.Scope = scopes
	}
}
