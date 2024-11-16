package pay

type Options struct {
	ServerAccountFile string
}

func WithServerAccountFile(file string) Option {
	return func(o *Options) {
		o.ServerAccountFile = file
	}
}
