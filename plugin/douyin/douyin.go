package douyin

import (
	"context"
	"fmt"

	"github.com/lsjhtang/kid/config"
	"github.com/lsjhtang/kid/injector"
	"github.com/lsjhtang/kid/utils"
)

const (
	douyinUrl = "https://open.douyin.com/oauth"
)

type Douyin struct {
	options     Options
	pay         *Payment
	app         *App
	miniProgram *App
}

type Option func(*Options)

func (douyin *Douyin) Provide(ctx context.Context) interface{} {
	var confName string
	if name, ok := ctx.Value(injector.NameKey{}).(string); ok && len(name) > 0 {
		confName = "." + name
	}
	confPrefix := fmt.Sprintf("douyin%s", confName)
	if !config.Exist(confPrefix) {
		panic(fmt.Sprintf("config.yaml file not found configuration item [%s]", confPrefix))
	}
	return New(
		WithAppClientKey(config.GetString(utils.GetConfigurationItem(confPrefix, "app.clientKey"))),
		WithAppClientSecret(config.GetString(utils.GetConfigurationItem(confPrefix, "app.clientSecret"))),
		WithMiniProgramClientKey(config.GetString(utils.GetConfigurationItem(confPrefix, "miniProgram.clientKey"))),
		WithMiniProgramClientSecret(config.GetString(utils.GetConfigurationItem(confPrefix, "miniProgram.clientSecret"))),
		WithPayAppId(config.GetString(utils.GetConfigurationItem(confPrefix, "pay.appId"))),
		WithPayDisableMsg(config.GetInt(utils.GetConfigurationItem(confPrefix, "pay.disableMsg"))),
		WithPayNotifyUrl(config.GetString(utils.GetConfigurationItem(confPrefix, "pay.notifyUrl"))),
		WithPayMsgPage(config.GetString(utils.GetConfigurationItem(confPrefix, "pay.msgPage"))),
		WithPaySalt(config.GetString(utils.GetConfigurationItem(confPrefix, "pay.salt"))),
		WithPayStoreUid(config.GetString(utils.GetConfigurationItem(confPrefix, "pay.storeUid"))),
		WithPayThirdpartyId(config.GetString(utils.GetConfigurationItem(confPrefix, "pay.thirdpartyId"))),
		WithPayToken(config.GetString(utils.GetConfigurationItem(confPrefix, "pay.token"))),
	)
}

// Payment 返回支付实例
func (douyin *Douyin) Payment() *Payment {
	if douyin.pay == nil {
		douyin.pay = NewPayment(douyin.options.PayOptions)
	}
	return douyin.pay
}

// App 返回APP实例
func (douyin *Douyin) App() *App {
	if douyin.app == nil {
		douyin.app = NewApp(douyin.options.AppOptions)
	}
	return douyin.app
}

// MiniProgram 返回小程序实例
func (douyin *Douyin) MiniProgram() *App {
	if douyin.miniProgram == nil {
		douyin.app = NewApp(douyin.options.MiniProgramOptions)
	}
	return douyin.app
}

func New(opts ...Option) *Douyin {
	options := Options{
		AppOptions:         &AppOptions{},
		MiniProgramOptions: &AppOptions{},
		PayOptions:         &PayOptions{},
	}
	for _, opt := range opts {
		opt(&options)
	}
	return &Douyin{
		options: options,
	}
}
