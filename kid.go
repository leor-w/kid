package kid

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/leor-w/kid/config"
	"github.com/leor-w/kid/container"
	"github.com/leor-w/kid/logger"
	"github.com/leor-w/kid/middleware"
	"github.com/leor-w/kid/plugin"
	"github.com/leor-w/kid/plugin/logrus"
)

var (
	loadConfigFlag = 0x1
	loadLoggerFlag = 0x2
)

const CommandLogo = `
        ██╗  ██╗  ██╗  ███████╗
        ██║ ██╔╝  ██║  ██╔═══██╗
        █████╔╝   ██║  ██║   ██║
        ██╔═██╗   ██║  ██║   ██║
        ██║  ██╗  ██║  ███████╔╝
        ╚═╝  ╚═╝  ╚═╝  ╚══════╝`

type Kid struct {
	loadFlag int
	*gin.Engine
	RouterGroup
	Logger       logger.Logger
	IocContainer *container.Container
	Options      *Options
}

var _defaultKid *Kid

type Option func(*Options)

func Registry(plugin plugin.Plugin, opts ...container.Option) *Kid {
	return _defaultKid.Registry(plugin, opts...)
}

// Registry 将插件注入容器
func (kid *Kid) Registry(plugin plugin.Plugin, opts ...container.Option) *Kid {
	if err := kid.IocContainer.Provide(plugin, opts...); err != nil {
		logger.Fatal(err.Error())
	}
	return kid
}

func (kid *Kid) Get(plugin plugin.Plugin, name ...string) (interface{}, error) {
	return kid.IocContainer.Get(plugin, name...)
}

func (kid *Kid) PopulateOne(plugin plugin.Plugin, opts ...container.Option) error {
	complete, err := kid.IocContainer.PopulateSingle(plugin, opts...)
	if err != nil || !complete {
		return fmt.Errorf("PopulateOne: 失败: %w", err)
	}
	return nil
}

func (kid *Kid) Populate() error {
	if kid.loadFlag < loadLoggerFlag|loadConfigFlag {
		panic("kid.Launch: failed: you must be call loadLogger or loadConfig")
	}
	return kid.IocContainer.Populate()
}

func (kid *Kid) Init() error {
	if err := kid.CommandLine(); err != nil {
		return err
	}
	// 读取配置文件
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
	kid.Engine.Use(middleware.Logger())
	kid.Engine.Use(middleware.Recovery)
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
	logger.Init(logrus.Default(""))
	kid.loadFlag |= loadLoggerFlag
	logger.Infof("日志服务准备就绪...")
	return nil
}

func (kid *Kid) Launch(hosts ...string) {
	host := ":8080"
	if len(hosts) > 0 {
		host = hosts[0]
	}
	PrintLogo()
	logger.Infof("HTTP 服务已启动: %s", host)
	if err := kid.Run(host); err != nil {
		logger.Error(err.Error())
	}
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
	_defaultKid = &Kid{
		IocContainer: container.New(),
		Options:      opt,
	}
	if err := _defaultKid.Init(); err != nil {
		panic(err.Error())
	}
	return _defaultKid
}

func PrintLogo() {
	logger.Info(CommandLogo)
}
