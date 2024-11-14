package tiktok

type Options struct {
	ClientKey    string
	ClientSecret string
	RedirectURL  string
	Scope        []string
}

func WithClientKey(clientKey string) Option {
	return func(o *Options) {
		o.ClientKey = clientKey
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
		case "basic":
			scopes = append(scopes, UserScopeBasic)
		case "artist_read":
			scopes = append(scopes, UserScopeArtistRead)
		case "artist_update":
			scopes = append(scopes, UserScopeArtistUpdate)
		case "profile":
			scopes = append(scopes, UserScopeUserProfile)
		case "stats":
			scopes = append(scopes, UserScopeUserStats)
		}
	}
	return func(o *Options) {
		o.Scope = scopes
	}
}
