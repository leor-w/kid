package rbac

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/leor-w/kid/config"
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
	return nil
}

func (ctrl *Controller) Provide() interface{} {
	if !config.Exist("rbac") {
		panic("not found [rbac] in config")
	}
	defaultRbac = New(
		WithModelConf(config.GetString("rbac.model")),
		WithEnableLog(config.GetBool("rbac.enableLog")),
		WithUser(config.GetString("rbac.db.user")),
		WithPwd(config.GetString("rbac.db.pwd")),
		WithHost(config.GetString("rbac.db.host")),
		WithPort(config.GetInt("rbac.db.port")),
		WithDbName(config.GetString("rbac.db.name")),
		WithTable(config.GetString("rbac.db.table")),
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
