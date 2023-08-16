package alipay

type Options struct {
	AppId             string
	PrivateKey        string
	AppPublicCert     string
	AppPublicCertFile string
	AliRootCert       string
	AliRootCertFile   string
	AliPublicCert     string
	AliPublicCertFile string
	IsProduct         bool
}

func WithAppId(appId string) Option {
	return func(o *Options) {
		o.AppId = appId
	}
}

func WithPrivateKey(privateKey string) Option {
	return func(o *Options) {
		o.PrivateKey = privateKey
	}
}

func WithAppPublicCert(publicCert string) Option {
	return func(o *Options) {
		o.AppPublicCert = publicCert
	}
}

func WithAppPublicCertFile(publicCertFile string) Option {
	return func(o *Options) {
		o.AppPublicCertFile = publicCertFile
	}
}

func WithAliRootCert(rootCert string) Option {
	return func(o *Options) {
		o.AliRootCert = rootCert
	}
}

func WithAliRootCertFile(rootCertFile string) Option {
	return func(o *Options) {
		o.AliRootCertFile = rootCertFile
	}
}

func WithAliPublicCert(publicCert string) Option {
	return func(o *Options) {
		o.AliPublicCert = publicCert
	}
}

func WithAliPublicCertFile(publicCertFile string) Option {
	return func(o *Options) {
		o.AliPublicCertFile = publicCertFile
	}
}

func WithIsProduct(isProduct bool) Option {
	return func(o *Options) {
		o.IsProduct = isProduct
	}
}
