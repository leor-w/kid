package simple

import "time"

type Option func(*Options)

type Options struct {
	// path 日志文件存储路径
	path string
	// logName 日志文件名称
	logName string
	// link 日志文件软连接到的位置
	link string
	// rotate 日志分割时间
	rotate time.Duration
	// maxAge 分割日志保存最长时间
	maxAge time.Duration
}

func WithPath(path string) Option {
	return func(o *Options) {
		o.path = path
	}
}

func WithLogName(logName string) Option {
	return func(o *Options) {
		o.logName = logName
	}
}

func WithLink(link string) Option {
	return func(o *Options) {
		o.link = link
	}
}

func WithRotate(rotate time.Duration) Option {
	return func(o *Options) {
		o.rotate = rotate
	}
}

func WithMaxAge(maxAge time.Duration) Option {
	return func(o *Options) {
		o.maxAge = maxAge
	}
}
