package kid

import "github.com/lsjhtang/kid/database/mysql"

type Options struct {
	Name            string
	AutoMigrate     bool
	Configs         []string
	RunMode         string
	DatabaseMigrate mysql.AutoMigrate
}

func WithAppName(name string) Option {
	return func(o *Options) {
		o.Name = name
	}
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

func WithDbMigrate(migrate mysql.AutoMigrate) Option {
	return func(o *Options) {
		o.DatabaseMigrate = migrate
	}
}
