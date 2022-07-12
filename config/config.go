package config

import (
	"errors"
	"github.com/leor-w/kid/config/local"
	"github.com/leor-w/utils"
	"github.com/spf13/cast"
	"net/url"
	"time"
)

type Config interface {
	Getter
	Init() error
	Provider(provider Provider)
}

type Provider interface {
	Get(key string) interface{}
	OnWatch(func()) error
	ReadConfig() error
	Unmarshal(key string, receiver interface{}) error
}

type Getter interface {
	GetString(key string) string
	GetStringSlice(key string) []string
	GetStringMap(key string) map[string]string
	GetStringMapStringSlice(string) map[string][]string
	GetInt(key string) int
	GetIntSlice(key string) []int
	GetInt64(key string) int64
	GetBool(key string) bool
	GetFloat(key string) float64
	GetTime(key string) time.Time
	GetDuration(key string) time.Duration
}

type kidConfig struct {
	providers []Provider
	opts      *Options
}

type Option func(*Options)

type Options struct {
	Providers []string
}

func WithProviders(providers []string) Option {
	return func(o *Options) {
		o.Providers = providers
	}
}

func (conf *kidConfig) Init() error {
	for _, provider := range conf.opts.Providers {
		_, err := url.Parse(provider)
		if err != nil {
			return errors.New("remote configuration not currently supported")
		}
		dir, name, ext := utils.ParsePath(provider)
		conf.Provider(local.New(local.WithConfigPath(dir),
			local.WithConfigName(name),
			local.WithConfigType(ext)))
	}
	return nil
}

func (conf *kidConfig) Provider(provider Provider) {
	conf.providers = append(conf.providers, provider)
}

func (conf *kidConfig) get(key string) interface{} {
	var val interface{}
	for _, provider := range conf.providers {
		val = provider.Get(key)
	}
	return val
}

func (conf *kidConfig) GetString(key string) string {
	return cast.ToString(conf.get(key))
}

func (conf *kidConfig) GetStringSlice(key string) []string {
	return cast.ToStringSlice(conf.get(key))
}

func (conf *kidConfig) GetStringMap(key string) map[string]string {
	return cast.ToStringMapString(conf.get(key))
}

func (conf *kidConfig) GetStringMapStringSlice(key string) map[string][]string {
	return cast.ToStringMapStringSlice(conf.get(key))
}

func (conf *kidConfig) GetInt(key string) int {
	return cast.ToInt(conf.get(key))
}

func (conf *kidConfig) GetIntSlice(key string) []int {
	return cast.ToIntSlice(conf.get(key))
}

func (conf *kidConfig) GetInt64(key string) int64 {
	return cast.ToInt64(conf.get(key))
}

func (conf *kidConfig) GetBool(key string) bool {
	return cast.ToBool(conf.get(key))
}

func (conf *kidConfig) GetFloat(key string) float64 {
	return cast.ToFloat64(conf.get(key))
}

func (conf *kidConfig) GetTime(key string) time.Time {
	return cast.ToTime(conf.get(key))
}

func (conf *kidConfig) GetDuration(key string) time.Duration {
	return cast.ToDuration(conf.get(key))
}

var (
	defaultConfig Config
)

func New(opts ...Option) Config {
	o := &Options{}
	for _, opt := range opts {
		opt(o)
	}
	conf := &kidConfig{
		opts: o,
	}
	return conf
}

func Default() Config {
	defaultConfig = &kidConfig{}
	return defaultConfig
}

func SetProvider(provider Provider) {
	defaultConfig.Provider(provider)
}

func GetString(key string) string {
	return defaultConfig.GetString(key)
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

func GetIntSlice(key string) []int {
	return defaultConfig.GetIntSlice(key)
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
