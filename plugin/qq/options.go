package qq

type Options struct {
	AppId       string
	AppKey      string
	RedirectUri string
	Display     string
	Product     bool
}

type Option func(*Options)

func WithAppId(appId string) Option {
	return func(o *Options) {
		o.AppId = appId
	}
}

func WithAppKey(appKey string) Option {
	return func(o *Options) {
		o.AppKey = appKey
	}
}

func WithProduct(product bool) Option {
	return func(o *Options) {
		o.Product = product
	}
}

func WithRedirectUri(redirectUri string) Option {
	return func(o *Options) {
		o.RedirectUri = redirectUri
	}
}

func WithDisplay(display string) Option {
	return func(o *Options) {
		o.Display = display
	}
}
