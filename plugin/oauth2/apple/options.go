package apple

import "fmt"

type Options struct {
	ClientId      string
	KeyId         string
	KeySecret     string
	KeySecretFile string
	TeamId        string
	RedirectURL   string
	Scopes        []string
}

func WithClientID(clientID string) Option {
	return func(o *Options) {
		o.ClientId = clientID
	}
}

func WithKeyId(keyId string) Option {
	return func(o *Options) {
		o.KeyId = keyId
	}
}

func WithKeySecret(clientSecret string) Option {
	return func(o *Options) {
		o.KeySecret = clientSecret
	}
}

func WithKeySecretFile(clientSecretFile string) Option {
	return func(o *Options) {
		o.KeySecretFile = clientSecretFile
	}
}

func WithTeamId(teamId string) Option {
	return func(o *Options) {
		o.TeamId = teamId
	}
}

func WithRedirectURL(redirectURL string) Option {
	return func(o *Options) {
		o.RedirectURL = redirectURL
	}
}

func WithScope(scope ...string) Option {
	scopes := make([]string, 0)
	for _, s := range scope {
		switch s {
		case ScopeOpenId, ScopeEmail, ScopeName:
			scopes = append(scopes, s)
		default:
			panic(fmt.Sprintf("Apple OAuth2: 无效的 Scopes: %s", s))
		}
	}
	return func(o *Options) {
		o.Scopes = scopes
	}
}
