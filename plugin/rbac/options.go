package rbac

type Options struct {
	model     string
	enableLog bool
	DbOptions
}

type DbOptions struct {
	user   string
	pwd    string
	host   string
	port   int
	dbname string
	table  string
}

func WithModelConf(modelConf string) Option {
	return func(o *Options) {
		o.model = modelConf
	}
}

func WithEnableLog(enableLog bool) Option {
	return func(o *Options) {
		o.enableLog = enableLog
	}
}

func WithUser(user string) Option {
	return func(o *Options) {
		o.user = user
	}
}

func WithPwd(pwd string) Option {
	return func(o *Options) {
		o.pwd = pwd
	}
}

func WithHost(host string) Option {
	return func(o *Options) {
		o.host = host
	}
}

func WithPort(port int) Option {
	return func(o *Options) {
		o.port = port
	}
}

func WithDbName(dbname string) Option {
	return func(o *Options) {
		o.dbname = dbname
	}
}

func WithTable(table string) Option {
	return func(o *Options) {
		o.table = table
	}
}
