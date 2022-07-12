package remote

type Options struct {
	provider  string
	path      string
	endpoints []*Endpoint
	confType  string
	gpgKey    string
}

type Endpoint struct {
	Url    string
	Path   string
	Secret string
}

func WithProvider(provider string) Option {
	return func(o *Options) {
		o.provider = provider
	}
}

func WithPath(path string) Option {
	return func(o *Options) {
		o.path = path
	}
}

func WithEndpoints(endpoints []*Endpoint) Option {
	return func(o *Options) {
		o.endpoints = endpoints
	}
}

func WithConfigType(confType string) Option {
	return func(o *Options) {
		o.confType = confType
	}
}

func WithGpgKey(gpgKey string) Option {
	return func(o *Options) {
		o.gpgKey = gpgKey
	}
}
