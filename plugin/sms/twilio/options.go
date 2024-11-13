package twilio

type Options struct {
	Analog     bool
	AccountSid string
	AuthToken  string
	ServiceSid string
	Edge       string
	Region     string
}

func WithAnalog(analog bool) Option {
	return func(o *Options) {
		o.Analog = analog
	}
}

func WithAccountSid(accountSid string) Option {
	return func(o *Options) {
		o.AccountSid = accountSid
	}
}

func WithAuthToken(authToken string) Option {
	return func(o *Options) {
		o.AuthToken = authToken
	}
}

func WithEdge(edge string) Option {
	return func(o *Options) {
		o.Edge = edge
	}
}

func WithRegion(region string) Option {
	return func(o *Options) {
		o.Region = region
	}
}

func WithServiceSid(serviceSid string) Option {
	return func(o *Options) {
		o.ServiceSid = serviceSid
	}
}
