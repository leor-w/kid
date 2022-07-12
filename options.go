package kid

type Options struct {
	Configs []string
}

func WithConfigs(configs ...string) Option {
	return func(o *Options) {
		o.Configs = configs
	}
}
