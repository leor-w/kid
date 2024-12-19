package apple

type Options struct {
	IsProduct      bool
	KeyId          string
	IssuerId       string
	Bid            string
	PrivateKey     string
	PrivateKeyFile string
	SharedSecret   string
}

func WithSandbox(isProduct bool) Option {
	return func(o *Options) {
		o.IsProduct = isProduct
	}
}

func WithKeyId(keyId string) Option {
	return func(o *Options) {
		o.KeyId = keyId
	}
}

func WithIssuerId(issuerId string) Option {
	return func(o *Options) {
		o.IssuerId = issuerId
	}
}

func WithBid(bid string) Option {
	return func(o *Options) {
		o.Bid = bid
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

func WithShareSecret(shareSecret string) Option {
	return func(o *Options) {
		o.SharedSecret = shareSecret
	}
}
