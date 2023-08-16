package douyin

type Options struct {
	AppOptions         *AppOptions
	MiniProgramOptions *AppOptions
	PayOptions         *PayOptions
}

type AppOptions struct {
	ClientKey    string
	ClientSecret string
}

func WithAppClientKey(clientKey string) Option {
	return func(o *Options) {
		o.AppOptions.ClientKey = clientKey
	}
}

func WithAppClientSecret(clientSecret string) Option {
	return func(o *Options) {
		o.AppOptions.ClientSecret = clientSecret
	}
}

func WithMiniProgramClientKey(clientKey string) Option {
	return func(o *Options) {
		o.MiniProgramOptions.ClientKey = clientKey
	}
}

func WithMiniProgramClientSecret(clientSecret string) Option {
	return func(o *Options) {
		o.MiniProgramOptions.ClientSecret = clientSecret
	}
}

type PayOptions struct {
	AppId        string
	Token        string
	Salt         string
	ThirdpartyId string
	StoreUid     string
	NotifyUrl    string
	DisableMsg   int
	MsgPage      string
	IsProduct    bool
}

func WithPayAppId(appId string) Option {
	return func(o *Options) {
		o.PayOptions.AppId = appId
	}
}

func WithPayToken(token string) Option {
	return func(o *Options) {
		o.PayOptions.Token = token
	}
}

func WithPaySalt(salt string) Option {
	return func(o *Options) {
		o.PayOptions.Salt = salt
	}
}

func WithPayThirdpartyId(thirdpartyId string) Option {
	return func(o *Options) {
		o.PayOptions.ThirdpartyId = thirdpartyId
	}
}

func WithPayStoreUid(storeUid string) Option {
	return func(o *Options) {
		o.PayOptions.StoreUid = storeUid
	}
}

func WithPayNotifyUrl(notifyUrl string) Option {
	return func(o *Options) {
		o.PayOptions.NotifyUrl = notifyUrl
	}
}

func WithPayDisableMsg(disableMsg int) Option {
	return func(o *Options) {
		o.PayOptions.DisableMsg = disableMsg
	}
}

func WithPayMsgPage(msgPage string) Option {
	return func(o *Options) {
		o.PayOptions.MsgPage = msgPage
	}
}

func WithIsProduct(isProduct bool) Option {
	return func(o *Options) {
		o.PayOptions.IsProduct = isProduct
	}
}
