package elasticsearch

import "time"

type Options struct {
	logLevel   string
	esAddress  []string
	esUser     string
	esPassword string
	cmd        string
	indexName  IndexName
	health     time.Duration
}

type IndexName func() string

func WithLogLevel(logLevel string) Option {
	return func(o *Options) {
		o.logLevel = logLevel
	}
}

func WithEsAddress(url ...string) Option {
	return func(o *Options) {
		o.esAddress = url
	}
}

func WithEsUser(user string) Option {
	return func(o *Options) {
		o.esUser = user
	}
}

func WithEsPassword(password string) Option {
	return func(o *Options) {
		o.esPassword = password
	}
}

func WithCmd(cmd string) Option {
	return func(o *Options) {
		o.cmd = cmd
	}
}

func WithIndexName(indexName IndexName) Option {
	return func(o *Options) {
		o.indexName = indexName
	}
}

func WithHealth(health time.Duration) Option {
	return func(o *Options) {
		o.health = health
	}
}
