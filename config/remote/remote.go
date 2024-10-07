package remote

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/spf13/viper"
)

type viperRemote struct {
	options *Options
	*viper.Viper
}

type Option func(*Options)

func (remote *viperRemote) Get(key string) interface{} {
	return remote.Viper.Get(key)
}

func (remote *viperRemote) ReadConfig() error {
	if err := remote.ReadRemoteConfig(); err != nil {
		return fmt.Errorf("config.yaml.remote.ReadConfig: read config.yaml failed: %w", err)
	}
	return nil
}

func (remote *viperRemote) Unmarshal(key string, receiver interface{}) error {
	iv := reflect.ValueOf(receiver)
	if iv.Kind() != reflect.Ptr || iv.IsNil() {
		return errors.New("config.yaml.remote.Unmarshal: unmarshal receiver mast not be a nil-pointer")
	}
	return remote.Viper.UnmarshalKey(key, receiver)
}

func (remote *viperRemote) OnWatch(watching func()) error {
	if err := remote.WatchRemoteConfigOnChannel(); err != nil {
		return err
	}
	go func() {
		for {
			time.Sleep(5 * time.Second)
			if err := remote.WatchRemoteConfig(); err != nil {
				continue
			}
			watching()
		}
	}()
	return nil
}

func (remote *viperRemote) Exist(key string) bool {
	return remote.IsSet(key)
}

func (remote *viperRemote) init() error {
	isSecret := len(remote.options.gpgKey) == 0
	for _, endpoint := range remote.options.endpoints {
		if isSecret {
			if err := remote.AddRemoteProvider(remote.options.provider, endpoint.Url, endpoint.Path); err != nil {
				return err
			}
		} else {
			if err := remote.AddSecureRemoteProvider(remote.options.provider, endpoint.Url, endpoint.Path, endpoint.Secret); err != nil {
				return err
			}
		}
	}
	remote.SetConfigType(remote.options.confType)
	if err := remote.ReadRemoteConfig(); err != nil {
		return err
	}
	return nil
}

func New(opts ...Option) *viperRemote {
	conf := viperRemote{
		options: &Options{},
		Viper:   viper.New(),
	}
	for _, opt := range opts {
		opt(conf.options)
	}
	if err := conf.init(); err != nil {
		panic(err.Error())
	}
	return &conf
}
