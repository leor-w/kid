package alipay

import (
	"context"
	"fmt"

	"github.com/lsjhtang/kid/injector"
	"github.com/smartwalle/alipay/v3"

	"github.com/lsjhtang/kid/config"
	"github.com/lsjhtang/kid/logger"
	"github.com/lsjhtang/kid/utils"
)

type Alipay struct {
	*alipay.Client
	options *Options
}

const (
	ProductCodeQuickMsecurityPay   = "QUICK_MSECURITY_PAY"    // 无线快捷支付
	ProductCodeCyclePayAuth        = "CYCLE_PAY_AUTH"         // 周期扣款产品
	ProductCodeFastInstantTradePay = "FAST_INSTANT_TRADE_PAY" // pc 网页端支付的支付方式
)

const (
	PayStatusTradeSuccess  = "TRADE_SUCCESS"  // 交易成功
	PayStatusTradeClose    = "TRADE_CLOSED"   // 交易关闭
	PayStatusTradeFinished = "TRADE_FINISHED" // 交易完成(不可退款)
)

type AlipaySDK string

const (
	AlipaySDKAPP  AlipaySDK = "APP" // app 支付
	AlipaySDKPage AlipaySDK = "PC"  // pc 网页端支付
	AliPaySDKWap  AlipaySDK = "WAP" // 手机网页端支付
)

func (pay *Alipay) Provide(ctx context.Context) interface{} {
	var confName string
	name, ok := ctx.Value(injector.NameKey{}).(string)
	if ok && len(name) > 0 {
		confName = "." + name
	}
	confPrefix := fmt.Sprintf("alipay%s", confName)
	if !config.Exist(confPrefix) {
		panic(fmt.Sprintf("config.yaml file not found configuration item [%s]", confPrefix))
	}
	return New(
		WithAppId(config.GetString(utils.GetConfigurationItem(confPrefix, "appId"))),
		WithPrivateKey(config.GetString(utils.GetConfigurationItem(confPrefix, "privateKey"))),
		WithAppPublicCert(config.GetString(utils.GetConfigurationItem(confPrefix, "appPublicCert"))),
		WithAppPublicCertFile(config.GetString(utils.GetConfigurationItem(confPrefix, "appPublicCertFile"))),
		WithAliRootCert(config.GetString(utils.GetConfigurationItem(confPrefix, "aliRootCert"))),
		WithAliRootCertFile(config.GetString(utils.GetConfigurationItem(confPrefix, "aliRootCertFile"))),
		WithAliPublicCert(config.GetString(utils.GetConfigurationItem(confPrefix, "aliPublicCert"))),
		WithAliPublicCertFile(config.GetString(utils.GetConfigurationItem(confPrefix, "aliPublicCertFile"))),
		WithIsProduct(config.GetBool(utils.GetConfigurationItem(confPrefix, "isProduct"))),
	)
}

type Option func(*Options)

func New(opts ...Option) *Alipay {
	var options Options
	for i := range opts {
		opts[i](&options)
	}
	cli, err := alipay.New(options.AppId, options.PrivateKey, options.IsProduct)
	if err != nil {
		logger.Errorf("初始化支付宝支付错误: %s", err.Error())
		return nil
	}
	pay := &Alipay{
		Client:  cli,
		options: &options,
	}
	if err = pay.Init(); err != nil {
		logger.Errorf("初始化支付宝支付错误: %s", err.Error())
		return nil
	}
	return pay
}

func (pay *Alipay) Init() error {
	if len(pay.options.AppPublicCert) == 0 {
		if err := pay.LoadAppCertPublicKeyFromFile(pay.options.AppPublicCertFile); err != nil {
			return err
		}
	} else {
		if err := pay.LoadAppCertPublicKey(pay.options.AppPublicCert); err != nil {
			return err
		}
	}
	if len(pay.options.AliRootCert) == 0 {
		if err := pay.LoadAliPayRootCertFromFile(pay.options.AliRootCertFile); err != nil {
			return err
		}
	} else {
		if err := pay.LoadAliPayRootCert(pay.options.AliRootCert); err != nil {
			return err
		}
	}
	if len(pay.options.AliPublicCert) == 0 {
		if err := pay.LoadAlipayCertPublicKeyFromFile(pay.options.AliPublicCertFile); err != nil {
			return err
		}
	} else {
		if err := pay.LoadAlipayCertPublicKey(pay.options.AliPublicCert); err != nil {
			return err
		}
	}
	return nil
}
