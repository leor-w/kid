package facebook

import "fmt"

type Options struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scope        []string
	Fields       []string
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
	var s []string
	for _, scope := range scopes {
		switch scope {
		case ScopeEmail:
			s = append(s, ScopeEmail)
		case ScopePublicProfile:
			s = append(s, ScopePublicProfile)
		default:
			panic(fmt.Sprintf("facebook oauth2: 未知的 scope 类型 [%s]", scope))
		}
	}
	return func(o *Options) {
		o.Scope = scopes
	}
}
