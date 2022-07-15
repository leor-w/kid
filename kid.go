package kid

import (
	"github.com/gin-gonic/gin"
	"github.com/leor-w/kid/config"
	"github.com/leor-w/kid/container"
	"github.com/leor-w/kid/logger"
	"github.com/leor-w/kid/plugin"
	"github.com/leor-w/kid/plugin/logrus"
)

var (
	loadConfigFlag = 0x1
	loadLoggerFlag = 0x2
)

type Kid struct {
	loadFlag int
	*gin.Engine
	RouterGroup
	Logger       logger.Logger
	iocContainer container.Container
	Options      *Options
}

type Option func(*Options)

// Registry 将插件注入容器
func (kid *Kid) Registry(plugin plugin.Plugin, opts ...container.Option) *Kid {
	if err := kid.iocContainer.Provide(plugin.Provide(), opts...); err != nil {
		logger.Fatal(err.Error())
	}
	return kid
}

func (kid *Kid) loadConfig() error {
	_ = config.New(config.WithProviders(kid.Options.Configs))
	kid.loadFlag |= loadConfigFlag
	return nil
}

func (kid *Kid) loadLogger() error {
	logger.Init(logrus.Default())
	kid.loadFlag |= loadLoggerFlag
	return nil
}

func (kid *Kid) Launch(hosts ...string) {
	host := ":8080"
	if len(hosts) > 0 {
		host = hosts[0]
	}
	if kid.loadFlag < loadLoggerFlag|loadConfigFlag {
		panic("kid.Launch: failed: you must be call loadLogger or loadConfig")
	}
	if err := kid.iocContainer.Populate(); err != nil {
		logger.Fatalf("kid.Launch: failed: %s", err.Error())
	}
	_ = kid.Run(host)
}

func (kid *Kid) NoRoute(handleFunc HandleFunc) {
	kid.Engine.NoRoute(convertHandleFunc(handleFunc))
}

func (kid *Kid) NoMethod(handleFunc HandleFunc) {
	kid.Engine.NoMethod(convertHandleFunc(handleFunc))
}

func New(opts ...Option) *Kid {
	opt := &Options{}
	for _, o := range opts {
		o(opt)
	}
	engine := gin.New()
	kid := &Kid{
		Engine:       engine,
		RouterGroup:  RouterGroup{&engine.RouterGroup},
		iocContainer: container.New(),
		Options:      opt,
	}
	if err := kid.loadLogger(); err != nil {
		panic(err.Error())
	}
	if err := kid.loadConfig(); err != nil {
		panic(err.Error())
	}
	return kid
}
