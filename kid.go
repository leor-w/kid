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
	*RouterGroup
	Logger       logger.Logger
	iocContainer container.Container
	Options      *Options
}

type Option func(*Options)

// Registry 将插件注入容器
func (kid *Kid) Registry(plugin plugin.Plugin, opts ...container.Option) *Kid {
	if err := kid.iocContainer.Provide(plugin, opts...); err != nil {
		logger.Fatal(err.Error())
	}
	return kid
}

func (kid *Kid) loadConfig() error {
	conf := config.New(config.WithProviders(kid.Options.Configs))
	if err := conf.Init(); err != nil {
		return err
	}
	kid.loadFlag |= loadConfigFlag
	return nil
}

func (kid *Kid) loadLogger() error {
	//logger.Default()
	log := logrus.Default()
	logger.Init(log)
	kid.loadFlag |= loadLoggerFlag
	return nil
}

func (kid *Kid) Launch(hosts ...string) {
	host := ":8080"
	if len(hosts) > 0 {
		host = hosts[0]
	}
	if err := kid.iocContainer.Populate(); err != nil {
		logger.Fatalf("kid.Launch: failed: %s", err.Error())
	}
	kid.Run(host)
}

func New(opts ...Option) *Kid {
	opt := &Options{}
	for _, o := range opts {
		o(opt)
	}
	engine := gin.New()
	kid := &Kid{
		Engine:       engine,
		RouterGroup:  &RouterGroup{&engine.RouterGroup},
		iocContainer: container.New(),
		Options:      opt,
	}
	if err := kid.loadConfig(); err != nil {
		panic(err.Error())
	}
	return kid
}
