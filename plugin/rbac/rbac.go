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
	host := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf-8", opt.user, opt.pwd, opt.host, opt.port, opt.dbname)
	adapter, err := gormadapter.NewAdapter("mysql", host, true)
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
	defaultRbac = New(
		WithModelConf("./conf/rbac_models.conf"),
		WithEnableLog(config.GetBool("rbac.enableLog")),
		WithUser(config.GetString("rbac.db.user")),
		WithPwd(config.GetString("rbac.db.pwd")),
		WithHost(config.GetString("rbac.db.host")),
		WithPort(config.GetInt("rbac.db.port")),
		WithDbName(config.GetString("rbac.db.name")),
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
