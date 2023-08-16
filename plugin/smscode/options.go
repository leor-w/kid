package smscode

type AliOptions struct {
	AccessKeyId     string
	AccessKeySecret string
	Endpoint        string
	SignName        string
	TempCode        string
}

func WithAccessKeyId(accessKeyId string) AliOption {
	return func(o *AliOptions) {
		o.AccessKeyId = accessKeyId
	}
}

func WithAccessKeySecret(accessKeySecret string) AliOption {
	return func(o *AliOptions) {
		o.AccessKeySecret = accessKeySecret
	}
}

func WithEndpoint(endpoint string) AliOption {
	return func(o *AliOptions) {
		o.Endpoint = endpoint
	}
}

func WithSignName(signName string) AliOption {
	return func(o *AliOptions) {
		o.SignName = signName
	}
}

func WithTempCode(tempCode string) AliOption {
	return func(o *AliOptions) {
		o.TempCode = tempCode
	}
}
