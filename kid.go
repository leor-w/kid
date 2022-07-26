package kid

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/leor-w/kid/config"
	"github.com/leor-w/kid/container"
	"github.com/leor-w/kid/logger"
	"github.com/leor-w/kid/middleware/ginmiddleware"
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

func (kid *Kid) Get(plugin plugin.Plugin, name ...string) (interface{}, error) {
	return kid.iocContainer.Get(plugin, name...)
}

func (kid *Kid) Init() error {
	if err := kid.loadConfig(); err != nil {
		return err
	}
	if err := kid.loadLogger(); err != nil {
		return err
	}
	runMode := config.GetString(kid.Options.RunMode)
	if len(runMode) == 0 {
		runMode = gin.TestMode
	}
	gin.SetMode(runMode)
	kid.Engine = gin.New()
	kid.Engine.Use(ginmiddleware.Logger())
	kid.Engine.Use(ginmiddleware.Recovery)
	kid.RouterGroup = RouterGroup{&(kid.Engine.RouterGroup)}
	return nil
}

func (kid *Kid) loadConfig() error {
	_ = config.New(
		config.WithProviders(kid.Options.Configs),
		config.WithDefault(true),
	)
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
		panic(fmt.Sprintf("kid.Launch: failed: %s", err.Error()))
	}
	logger.Infof("Listener and serving HTTP on: %s", host)
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
	kid := &Kid{
		iocContainer: container.New(),
		Options:      opt,
	}
	if err := kid.Init(); err != nil {
		panic(err.Error())
	}
	return kid
}
