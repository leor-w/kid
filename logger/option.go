package logger

import (
	"context"
)

type Options struct {
	Level   Level
	LogName string
	Context context.Context
}

func WithLevel(level Level) Option {
	return func(o *Options) {
		o.Level = level
	}
}

func SetOption(key, value interface{}) Option {
	return func(o *Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, key, value)
	}
}
