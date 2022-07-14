package wechat

type Options struct {
	*MiniProgram
	*Pay
	appid string
}

type MiniProgram struct {
	secret string
}

type Pay struct {
	mchid     string
	key       string
	notifyUrl string
}

func WithAppid(appid string) Option {
	return func(o *Options) {
		o.appid = appid
	}
}

func WithMiniSecret(secret string) Option {
	return func(o *Options) {
		o.MiniProgram.secret = secret
	}
}

func WithPayMchid(mchid string) Option {
	return func(o *Options) {
		o.Pay.mchid = mchid
	}
}

func WithPayKey(key string) Option {
	return func(o *Options) {
		o.Pay.key = key
	}
}

func WithPayNotifyUrl(notifyUrl string) Option {
	return func(o *Options) {
		o.Pay.notifyUrl = notifyUrl
	}
}
