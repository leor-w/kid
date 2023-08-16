package mysql

import (
	"time"
)

type Options struct {
	User     string        // 数据库用户名
	Password string        // 数据库密码
	Host     string        // 数据库连接地址
	Port     int           // 数据库连接端口
	Db       string        // 数据库库名
	MaxIdle  int           // 最大空闲连接数
	MaxOpen  int           // 最大连接数
	MaxLife  time.Duration // 连接最大存活时间
	LogLevel int           // 日志级别

	CloseFKCheck bool // 关闭外键检查
}

func WithUser(user string) Option {
	return func(o *Options) {
		o.User = user
	}
}

func WithPassword(password string) Option {
	return func(o *Options) {
		o.Password = password
	}
}

func WithHost(host string) Option {
	return func(o *Options) {
		o.Host = host
	}
}

func WithPort(port int) Option {
	return func(o *Options) {
		o.Port = port
	}
}

func WithDb(dbName string) Option {
	return func(o *Options) {
		o.Db = dbName
	}
}

func WithMaxIdle(maxIdle int) Option {
	return func(o *Options) {
		o.MaxIdle = maxIdle
	}
}

func WithMaxOpen(maxOpen int) Option {
	return func(o *Options) {
		o.MaxOpen = maxOpen
	}
}

func WithMaxLife(maxLife time.Duration) Option {
	return func(o *Options) {
		o.MaxLife = maxLife
	}
}

func WithLogLevel(logLevel int) Option {
	return func(o *Options) {
		o.LogLevel = logLevel
	}
}

func WithCloseFKCheck(closeFKCheck bool) Option {
	return func(o *Options) {
		o.CloseFKCheck = closeFKCheck
	}
}
