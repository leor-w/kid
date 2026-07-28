package local

import (
	"errors"

	"github.com/fsnotify/fsnotify"
	"github.com/lsjhtang/kid/utils"
	"github.com/spf13/viper"
)

type viperLocal struct {
	options *Options
	*viper.Viper
}

type Option func(*Options)

func (local *viperLocal) Get(key string) interface{} {
	return local.Viper.Get(key)
}

func (local *viperLocal) ReadConfig() error {
	return local.Viper.ReadInConfig()
}

func (local *viperLocal) OnWatch(watching func()) error {
	local.Viper.WatchConfig()
	local.Viper.OnConfigChange(func(in fsnotify.Event) {
		watching()
	})
	return nil
}

func (local *viperLocal) Unmarshal(key string, receiver interface{}) error {
	if utils.IsNilPointer(receiver) {
		return errors.New("config.yaml.local.Unmarshal: unmarshal receiver mast not be a nil-pointer")
	}
	return local.Viper.UnmarshalKey(key, &receiver)
}

func (local *viperLocal) Exist(key string) bool {
	return local.Viper.IsSet(key)
}

func New(opts ...Option) *viperLocal {
	conf := &viperLocal{
		options: &Options{
			configName: "config.yaml",
			configPath: "./",
			configType: "yaml",
		},
		Viper: viper.New(),
	}
	for _, opt := range opts {
		opt(conf.options)
	}
	conf.SetConfigName(conf.options.configName)
	conf.SetConfigType(conf.options.configType)
	conf.AddConfigPath(conf.options.configPath)
	return conf
}
