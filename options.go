package kid

type Options struct {
	Configs []string
	RunMode string
}

func WithConfigs(configs ...string) Option {
	return func(o *Options) {
		if len(o.Configs) > 0 {
			o.Configs = append(o.Configs, configs...)
			return
		}
		o.Configs = configs
	}
}

func WithRunMode(runMode string) Option {
	return func(o *Options) {
		o.RunMode = runMode
	}
}
