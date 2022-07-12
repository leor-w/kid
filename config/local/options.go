package local

type Options struct {
	configName string
	configPath string
	configType string
}

func WithConfigName(configName string) Option {
	return func(o *Options) {
		o.configName = configName
	}
}

func WithConfigPath(configPath string) Option {
	return func(o *Options) {
		o.configPath = configPath
	}
}

func WithConfigType(confType string) Option {
	return func(o *Options) {
		o.configType = confType
	}
}
