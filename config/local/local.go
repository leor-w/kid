package local

import (
	"errors"

	"github.com/fsnotify/fsnotify"
	"github.com/leor-w/utils"
	"github.com/spf13/viper"
)

type viperLocal struct {
	options *Options
	conf    *viper.Viper
}

type Option func(*Options)

func (local *viperLocal) Get(key string) interface{} {
	return local.conf.Get(key)
}

func (local *viperLocal) ReadConfig() error {
	return local.conf.ReadInConfig()
}

func (local *viperLocal) OnWatch(watching func()) error {
	local.conf.WatchConfig()
	local.conf.OnConfigChange(func(in fsnotify.Event) {
		watching()
	})
	return nil
}

func (local *viperLocal) Unmarshal(key string, receiver interface{}) error {
	if utils.IsNilPointer(receiver) {
		return errors.New("config.local.Unmarshal: unmarshal receiver mast not be a nil-pointer")
	}
	return local.conf.UnmarshalKey(key, &receiver)
}

func New(opts ...Option) *viperLocal {
	conf := &viperLocal{
		options: &Options{
			configName: "config",
			configPath: "./",
			configType: "yaml",
		},
		conf: viper.New(),
	}
	for _, opt := range opts {
		opt(conf.options)
	}
	conf.conf.SetConfigName(conf.options.configName)
	conf.conf.SetConfigType(conf.options.configType)
	conf.conf.AddConfigPath(conf.options.configPath)
	return conf
}
