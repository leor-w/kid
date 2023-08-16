package kdn

type Options struct {
	EBusinessID string
	AppKey      string
	DataType    string
	RequestType string
	BaseURL     string
}

func WithEBusinessID(eBusinessID string) Option {
	return func(o *Options) {
		o.EBusinessID = eBusinessID
	}
}

func WithAppKey(appKey string) Option {
	return func(o *Options) {
		o.AppKey = appKey
	}
}

func WithDataType(dataType string) Option {
	return func(o *Options) {
		o.DataType = dataType
	}
}

func WithRequestType(requestType string) Option {
	return func(o *Options) {
		o.RequestType = requestType
	}
}

func WithBaseURL(baseURL string) Option {
	return func(o *Options) {
		o.BaseURL = baseURL
	}
}
