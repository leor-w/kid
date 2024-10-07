package apple

type Options struct {
	ClientId         string
	KeyId            string
	ClientSecret     string
	ClientSecretFile string
	TeamId           string
	RedirectURL      string
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

func WithClientSecret(clientSecret string) Option {
	return func(o *Options) {
		o.ClientSecret = clientSecret
	}
}

func WithClientSecretFile(clientSecretFile string) Option {
	return func(o *Options) {
		o.ClientSecretFile = clientSecretFile
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
