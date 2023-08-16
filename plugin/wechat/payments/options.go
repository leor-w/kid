package payments

type Options struct {
	Appid            string
	Mchid            string
	PrivateKey       string
	PrivateKeyFile   string
	Certificate      string
	CertificateFile  string
	MchCertSerialNum string
	Key              string
	PayNotifyUrl     string
	RefundNotifyUrl  string
}

func WithAppId(appId string) Option {
	return func(o *Options) {
		o.Appid = appId
	}
}

func WithMchId(mchId string) Option {
	return func(o *Options) {
		o.Mchid = mchId
	}
}

func WithPrivateKey(privateKey string) Option {
	return func(o *Options) {
		o.PrivateKey = privateKey
	}
}

func WithPrivateKeyFile(privateKeyFile string) Option {
	return func(o *Options) {
		o.PrivateKeyFile = privateKeyFile
	}
}

func WithCertificate(cert string) Option {
	return func(o *Options) {
		o.Certificate = cert
	}
}

func WithCertificateFile(certFile string) Option {
	return func(o *Options) {
		o.CertificateFile = certFile
	}
}

func WithCertSerialNum(certSerialNum string) Option {
	return func(o *Options) {
		o.MchCertSerialNum = certSerialNum
	}
}

func WithKey(key string) Option {
	return func(o *Options) {
		o.Key = key
	}
}

func WithNotifyUrl(notifyUrl string) Option {
	return func(o *Options) {
		o.PayNotifyUrl = notifyUrl
	}
}

func WithRefundNotifyUrl(refundNotifyUrl string) Option {
	return func(o *Options) {
		o.RefundNotifyUrl = refundNotifyUrl
	}
}
