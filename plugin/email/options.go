package email

type Options struct {
	Host string // 邮箱服务地址
	Port int    // 邮箱服务端口
	From string // 邮件发送使用地址
	Pwd  string // 邮件密码
}

type Option func(*Options)

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

func WithFrom(from string) Option {
	return func(o *Options) {
		o.From = from
	}
}

func WithPwd(pwd string) Option {
	return func(o *Options) {
		o.Pwd = pwd
	}
}
