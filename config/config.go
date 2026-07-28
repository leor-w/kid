package config

import (
	"time"

	"github.com/lsjhtang/kid/config/local"
	"github.com/lsjhtang/kid/utils"
	"github.com/spf13/cast"
)

type Config struct {
	providers []Provider
	options   *Options
}

type Provider interface {
	Get(key string) interface{}
	OnWatch(func()) error
	ReadConfig() error
	Unmarshal(key string, receiver interface{}) error
	Exist(key string) bool
}

type Option func(*Options)

type Options struct {
	Providers []string
	Default   bool
}

func WithProviders(providers []string) Option {
	return func(o *Options) {
		o.Providers = providers
	}
}

func WithDefault(isDefault bool) Option {
	return func(o *Options) {
		o.Default = isDefault
	}
}

func (conf *Config) Init() error {
	for _, provider := range conf.options.Providers {
		dir, name, ext := utils.ParsePath(provider)
		localProvider := local.New(local.WithConfigPath(dir),
			local.WithConfigName(name),
			local.WithConfigType(ext))
		if err := localProvider.ReadInConfig(); err != nil {
			return err
		}
		conf.Provider(localProvider)
	}
	if conf.options.Default {
		defaultConfig = conf
	}
	return nil
}

func (conf *Config) Provider(provider Provider) {
	conf.providers = append(conf.providers, provider)
}

func (conf *Config) find(key string) interface{} {
	var val interface{}
	for _, provider := range conf.providers {
		val = provider.Get(key)
		if len(cast.ToString(val)) > 0 {
			return val
		}
	}
	return val
}

func (conf *Config) Get(key string) interface{} {
	var val interface{}
	for _, provider := range conf.providers {
		val = provider.Get(key)
		if val != nil {
			return val
		}
	}
	return val
}

func (conf *Config) GetString(key string) string {
	return cast.ToString(conf.find(key))
}

func (conf *Config) GetStringSlice(key string) []string {
	return cast.ToStringSlice(conf.find(key))
}

func (conf *Config) GetStringMap(key string) map[string]string {
	return cast.ToStringMapString(conf.find(key))
}

func (conf *Config) GetStringMapStringSlice(key string) map[string][]string {
	return cast.ToStringMapStringSlice(conf.find(key))
}

func (conf *Config) GetInt(key string) int {
	return cast.ToInt(conf.find(key))
}

func (conf *Config) GetUint(key string) uint {
	return cast.ToUint(conf.find(key))
}

func (conf *Config) GetIntSlice(key string) []int {
	return cast.ToIntSlice(conf.find(key))
}

func (conf *Config) GetInt64(key string) int64 {
	return cast.ToInt64(conf.find(key))
}

func (conf *Config) GetBool(key string) bool {
	return cast.ToBool(conf.find(key))
}

func (conf *Config) GetFloat(key string) float64 {
	return cast.ToFloat64(conf.find(key))
}

func (conf *Config) GetTime(key string) time.Time {
	return cast.ToTime(conf.find(key))
}

func (conf *Config) GetDuration(key string) time.Duration {
	return cast.ToDuration(conf.find(key))
}

func (conf *Config) Exist(key string) bool {
	for _, provider := range conf.providers {
		if provider.Exist(key) {
			return true
		}
	}
	return false
}

var (
	defaultConfig *Config
)

func New(opts ...Option) *Config {
	options := &Options{}
	for _, opt := range opts {
		opt(options)
	}
	conf := &Config{
		options: options,
	}
	if err := conf.Init(); err != nil {
		panic(err.Error())
	}
	return conf
}

func Default() *Config {
	defaultConfig = New(WithProviders([]string{"./config.yaml.yaml"}))
	return defaultConfig
}

func SetProvider(provider Provider) {
	defaultConfig.Provider(provider)
}

func Get(key string) interface{} {
	return defaultConfig.Get(key)
}

func GetString(key string) string {
	return defaultConfig.GetString(key)
}

func DefaultString(key, defaultVal string) string {
	val := GetString(key)
	if len(val) == 0 {
		return defaultVal
	}
	return val
}

func GetStringSlice(key string) []string {
	return defaultConfig.GetStringSlice(key)
}

func GetStringMap(key string) map[string]string {
	return defaultConfig.GetStringMap(key)
}

func GetStringMapStringSlice(key string) map[string][]string {
	return defaultConfig.GetStringMapStringSlice(key)
}

func GetInt(key string) int {
	return defaultConfig.GetInt(key)
}

func GetUint(key string) uint {
	return defaultConfig.GetUint(key)
}

func GetIntSlice(key string) []int {
	return defaultConfig.GetIntSlice(key)
}

func GetUint32(key string) uint32 {
	return uint32(defaultConfig.GetUint(key))
}

func GetInt64(key string) int64 {
	return defaultConfig.GetInt64(key)
}

func GetBool(key string) bool {
	return defaultConfig.GetBool(key)
}

func GetFloat(key string) float64 {
	return defaultConfig.GetFloat(key)
}

func GetTime(key string) time.Time {
	return defaultConfig.GetTime(key)
}

func GetDuration(key string) time.Duration {
	return defaultConfig.GetDuration(key)
}

func Exist(key string) bool {
	return defaultConfig.Exist(key)
}
