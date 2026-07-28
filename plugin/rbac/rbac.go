package rbac

import (
	"context"
	"fmt"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/lsjhtang/kid/config"
	"github.com/lsjhtang/kid/injector"
	"github.com/lsjhtang/kid/utils"
)

type Controller struct {
	*casbin.Enforcer
	adapter *gormadapter.Adapter
	options *Options
}

type Option func(*Options)

func (ctrl *Controller) Init() error {
	opt := ctrl.options
	host := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", opt.user, opt.pwd, opt.host, opt.port, opt.dbname)
	adapter, err := gormadapter.NewAdapter("mysql", host, opt.dbname, opt.table, true)
	if err != nil {
		return err
	}
	enforcer, err := casbin.NewEnforcer(opt.model, adapter)
	if err != nil {
		return err
	}
	if err = enforcer.LoadPolicy(); err != nil {
		return err
	}
	enforcer.EnableLog(opt.enableLog)
	ctrl.adapter = adapter
	ctrl.Enforcer = enforcer
	return nil
}

func (ctrl *Controller) Provide(ctx context.Context) interface{} {
	var confName string
	name, ok := ctx.Value(injector.NameKey{}).(string)
	if ok && len(name) > 0 {
		confName = "." + name
	}
	confPrefix := fmt.Sprintf("rbac%s", confName)
	if !config.Exist(confPrefix) {
		panic(fmt.Sprintf("config.yaml file not found configuration item [%s]", confPrefix))
	}
	defaultRbac = New(
		WithModelConf(config.GetString(utils.GetConfigurationItem(confPrefix, "model"))),
		WithEnableLog(config.GetBool(utils.GetConfigurationItem(confPrefix, "enableLog"))),
		WithUser(config.GetString(utils.GetConfigurationItem(confPrefix, "db.user"))),
		WithPwd(config.GetString(utils.GetConfigurationItem(confPrefix, "db.pwd"))),
		WithHost(config.GetString(utils.GetConfigurationItem(confPrefix, "db.host"))),
		WithPort(config.GetInt(utils.GetConfigurationItem(confPrefix, "db.port"))),
		WithDbName(config.GetString(utils.GetConfigurationItem(confPrefix, "db.name"))),
		WithTable(config.GetString(utils.GetConfigurationItem(confPrefix, "db.table"))),
	)
	return defaultRbac
}

var (
	defaultRbac *Controller
)

func Default() *Controller {
	return defaultRbac
}

func New(opts ...Option) *Controller {
	var options Options
	for _, o := range opts {
		o(&options)
	}
	ctrl := &Controller{
		options: &options,
	}
	if err := ctrl.Init(); err != nil {
		panic(err.Error())
	}
	return ctrl
}
