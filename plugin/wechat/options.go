package wechat

import "github.com/lsjhtang/kid/plugin/wechat/payments"

type Options struct {
	MiniProgram             // 小程序配置
	OfficialAccount         // 公众号配置
	payments.Options        // 支付配置
	appid            string // 共有 appId 可被各自配置覆盖
}

type MiniProgram struct {
	appId  string
	secret string
}

type OfficialAccount struct {
	appId          string
	token          string
	secret         string
	encodingAesKey string
}

func WithAppId(appId string) Option {
	return func(o *Options) {
		o.appid = appId
	}
}

func WithMiniAppId(appId string) Option {
	return func(o *Options) {
		o.MiniProgram.appId = appId
	}
}

func WithMiniSecret(secret string) Option {
	return func(o *Options) {
		o.MiniProgram.secret = secret
	}
}

func WithPayAppId(appId string) Option {
	return func(o *Options) {
		o.Options.Appid = appId
	}
}

func WithPayMchid(mchid string) Option {
	return func(o *Options) {
		o.Mchid = mchid
	}
}

func WithPayKey(key string) Option {
	return func(o *Options) {
		o.Key = key
	}
}

func WithPayPrivateKey(privateKey string) Option {
	return func(o *Options) {
		o.PrivateKey = privateKey
	}
}

func WithPayPrivateKeyFile(privateKeyFile string) Option {
	return func(o *Options) {
		o.PrivateKeyFile = privateKeyFile
	}
}

func WithCertSerialNum(certSerialNum string) Option {
	return func(o *Options) {
		o.MchCertSerialNum = certSerialNum
	}
}

func WithPayNotifyUrl(notifyUrl string) Option {
	return func(o *Options) {
		o.PayNotifyUrl = notifyUrl
	}
}

func WithOfficialAppId(appId string) Option {
	return func(o *Options) {
		o.OfficialAccount.appId = appId
	}
}

func WithOfficialToken(token string) Option {
	return func(o *Options) {
		o.token = token
	}
}

func WithOfficialEncodingAesKey(key string) Option {
	return func(o *Options) {
		o.encodingAesKey = key
	}
}

func WithOfficialSecret(secret string) Option {
	return func(o *Options) {
		o.OfficialAccount.secret = secret
	}
}
