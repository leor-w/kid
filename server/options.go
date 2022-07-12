package server

import "crypto/tls"

type Options struct {
	// Address 为服务指定地址 ex: :8080
	Address string
	// TLSConfig 为服务指定 TLS 配置
	TLSConfig *tls.Config
}

func Address(addr string) Option {
	return func(o *Options) {
		o.Address = addr
	}
}

func TLSConfig(conf *tls.Config) Option {
	return func(o *Options) {
		o.TLSConfig = conf
	}
}
