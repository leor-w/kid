package rate_limit

type Options struct {
	RateLimit       int // 每秒允许的最大请求数
	BurstLimit      int // 最大突发请求数
	RateLimitWindow int // 速率限制窗口大小，单位秒
}

func WithRateLimit(rateLimit int) Option {
	return func(o *Options) {
		o.RateLimit = rateLimit
	}
}

func WithBurstLimit(burstLimit int) Option {
	return func(o *Options) {
		o.BurstLimit = burstLimit
	}
}

func WithRateLimitWindow(rateLimitWindow int) Option {
	return func(o *Options) {
		o.RateLimitWindow = rateLimitWindow
	}
}
