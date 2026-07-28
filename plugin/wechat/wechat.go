package wechat

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"

	"github.com/lsjhtang/kid/injector"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/miniprogram"
	miniConf "github.com/silenceper/wechat/v2/miniprogram/config"
	"github.com/silenceper/wechat/v2/officialaccount"
	oaConf "github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/silenceper/wechat/v2/pay"
	v2PayConf "github.com/silenceper/wechat/v2/pay/config"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/auth/verifiers"
	"github.com/wechatpay-apiv3/wechatpay-go/core/downloader"
	"github.com/wechatpay-apiv3/wechatpay-go/core/notify"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/services/refunddomestic"
	v3Utils "github.com/wechatpay-apiv3/wechatpay-go/utils"

	"github.com/lsjhtang/kid/config"
	"github.com/lsjhtang/kid/plugin/wechat/payments"
	"github.com/lsjhtang/kid/utils"
)

const (
	TradeSuccess    = "SUCCESS"    // 支付成功
	TradeRefund     = "REFUND"     // 转入退款
	TradeNotPay     = "NOTPAY"     // 未支付
	TradeClosed     = "CLOSED"     // 已关闭
	TradeRevoked    = "REVOKED"    // 已撤销
	TradeUserPaying = "USERPAYING" // 用户支付中
	TradePayError   = "PAYERROR"   // 支付失败
)

type Wechat struct {
	options *Options

	// 公众号及小程序相关
	wechat      *wechat.Wechat
	miniProgram *miniprogram.MiniProgram
	official    *officialaccount.OfficialAccount

	// API V2 版本 支付
	pay   *pay.Pay
	cache cache.Cache

	// API V3 版本 支付
	notify   *notify.Handler
	payments map[payments.WechatPaySDK]payments.Payment
	client   *core.Client
	refund   *refunddomestic.RefundsApiService
}

func (w *Wechat) Provide(ctx context.Context) interface{} {
	var confName string
	if name, ok := ctx.Value(injector.NameKey{}).(string); ok && len(name) > 0 {
		confName = "." + name
	}
	confPrefix := fmt.Sprintf("wechat%s", confName)
	if !config.Exist(confPrefix) {
		panic(fmt.Sprintf("config.yaml file not found configuration item [%s]", confPrefix))
	}
	return New(
		WithAppId(config.GetString(utils.GetConfigurationItem(confPrefix, "appid"))),
		WithMiniAppId(config.GetString(utils.GetConfigurationItem(confPrefix, "mini.appId"))),
		WithMiniSecret(config.GetString(utils.GetConfigurationItem(confPrefix, "mini.secret"))),
		WithPayAppId(config.GetString(utils.GetConfigurationItem(confPrefix, "pay.appId"))),
		WithPayPrivateKey(config.GetString(utils.GetConfigurationItem(confPrefix, "pay.privateKey"))),
		WithPayPrivateKeyFile(config.GetString(utils.GetConfigurationItem(confPrefix, "pay.privateKeyFile"))),
		WithPayMchid(config.GetString(utils.GetConfigurationItem(confPrefix, "pay.mchid"))),
		WithPayKey(config.GetString(utils.GetConfigurationItem(confPrefix, "pay.key"))),
		WithCertSerialNum(config.GetString(utils.GetConfigurationItem(confPrefix, "pay.certSerialNum"))),
		WithPayNotifyUrl(config.GetString(utils.GetConfigurationItem(confPrefix, "pay.payNotifyUrl"))),
		WithOfficialAppId(config.GetString(utils.GetConfigurationItem(confPrefix, "official.appId"))),
		WithOfficialToken(config.GetString(utils.GetConfigurationItem(confPrefix, "official.token"))),
		WithOfficialEncodingAesKey(config.GetString(utils.GetConfigurationItem(confPrefix, "official.encodingAESKey"))),
		WithOfficialSecret(config.GetString(utils.GetConfigurationItem(confPrefix, "official.secret"))),
	)
}

type Option func(*Options)

// ------------------ V2 版本 ------------------

func (w *Wechat) SetCache(cache cache.Cache) {
	w.cache = cache
	w.wechat.SetCache(cache)
}

// MiniProgram 微信小程序
func (w *Wechat) MiniProgram() *miniprogram.MiniProgram {
	if w.miniProgram == nil {
		appId := w.options.appid
		if len(w.options.MiniProgram.appId) > 0 {
			appId = w.options.MiniProgram.appId
		}
		w.miniProgram = w.wechat.GetMiniProgram(&miniConf.Config{
			AppID:     appId,
			AppSecret: w.options.MiniProgram.secret,
			Cache:     w.cache,
		})
	}
	return w.miniProgram
}

// OfficialAccount 公众号
func (w *Wechat) OfficialAccount() *officialaccount.OfficialAccount {
	if w.official == nil {
		appId := w.options.appid
		if len(w.options.OfficialAccount.appId) > 0 {
			appId = w.options.OfficialAccount.appId
		}
		w.official = w.wechat.GetOfficialAccount(&oaConf.Config{
			AppID:          appId,
			AppSecret:      w.options.OfficialAccount.secret,
			Token:          w.options.OfficialAccount.token,
			EncodingAESKey: w.options.OfficialAccount.encodingAesKey,
			Cache:          w.cache,
		})
	}
	return w.official
}

// Pay 微信支付 使用 v2 版本的API支付接口
func (w *Wechat) Pay() *pay.Pay {
	if w.pay == nil {
		appId := w.options.appid
		if len(w.options.Options.Appid) > 0 {
			appId = w.options.Options.Appid
		}
		w.pay = w.wechat.GetPay(&v2PayConf.Config{
			AppID:     appId,
			MchID:     w.options.Mchid,
			Key:       w.options.Key,
			NotifyURL: w.options.PayNotifyUrl,
		})
	}
	return w.pay
}

// ------------------ V2 版本 ------------------

// ------------------ V3 版本 ------------------

// Payment 微信支付 使用 V3 版本的支付 API 接口
func (w *Wechat) Payment(payType payments.WechatPaySDK) (payments.Payment, error) {
	if w.client == nil {
		var (
			mchPrivateKey *rsa.PrivateKey
			err           error
		)
		if len(w.options.PrivateKeyFile) > 0 {
			mchPrivateKey, err = v3Utils.LoadPrivateKeyWithPath(w.options.PrivateKeyFile)
		} else {
			err = errors.New("商户私钥为空, 无法初始化支付客户端")
		}
		if err != nil {
			return nil, err
		}
		ctx := context.Background()
		opts := []core.ClientOption{
			option.WithWechatPayAutoAuthCipher(w.options.Mchid, w.options.MchCertSerialNum, mchPrivateKey, w.options.Key),
		}
		client, err := core.NewClient(ctx, opts...)
		if err != nil {
			return nil, err
		}
		w.client = client
	}
	payment, exist := w.payments[payType]
	if !exist || payment == nil {
		var err error
		appId := w.options.appid
		if len(w.options.Options.Appid) > 0 {
			appId = w.options.Options.Appid
		}
		payment, err = payments.New(payType, w.client,
			payments.WithAppId(appId),
			payments.WithMchId(w.options.Mchid),
			payments.WithPrivateKey(w.options.PrivateKey),
			payments.WithPrivateKeyFile(w.options.PrivateKeyFile),
			payments.WithCertSerialNum(w.options.MchCertSerialNum),
			payments.WithKey(w.options.Key),
			payments.WithNotifyUrl(w.options.PayNotifyUrl),
			payments.WithRefundNotifyUrl(w.options.RefundNotifyUrl),
		)
		if err != nil {
			return nil, err
		}
		w.payments[payType] = payment
	}
	return payment, nil
}

// Notify 微信支付 V3 版本 API 回调通知处理器
func (w *Wechat) Notify() (*notify.Handler, error) {
	if w.client == nil {
		cli, err := w.newClient()
		if err != nil {
			return nil, err
		}
		w.client = cli
	}
	if w.notify == nil {
		ctx := context.Background()
		err := downloader.MgrInstance().RegisterDownloaderWithClient(ctx, w.client, w.options.Mchid, w.options.Key)
		if err != nil {
			return nil, err
		}
		certificateVisitor := downloader.MgrInstance().GetCertificateVisitor(w.options.Mchid)
		handler, err := notify.NewRSANotifyHandler(w.options.Key, verifiers.NewSHA256WithRSAVerifier(certificateVisitor))
		w.notify = handler
	}
	return w.notify, nil
}

func (w *Wechat) newClient() (*core.Client, error) {
	var (
		mchPrivateKey *rsa.PrivateKey
		err           error
	)
	if len(w.options.PrivateKeyFile) > 0 {
		mchPrivateKey, err = v3Utils.LoadPrivateKeyWithPath(w.options.PrivateKeyFile)
	} else {
		err = errors.New("商户私钥为空, 无法初始化支付客户端")
	}
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	opts := []core.ClientOption{
		option.WithWechatPayAutoAuthCipher(w.options.Mchid, w.options.MchCertSerialNum, mchPrivateKey, w.options.Key),
	}
	client, err := core.NewClient(ctx, opts...)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (w *Wechat) Refund() (*refunddomestic.RefundsApiService, error) {
	if w.client == nil {
		cli, err := w.newClient()
		if err != nil {
			return nil, err
		}
		w.client = cli
	}
	if w.refund == nil {
		w.refund = &refunddomestic.RefundsApiService{Client: w.client}
	}
	return w.refund, nil
}

// ------------------ V3 版本 ------------------

func New(opts ...Option) *Wechat {
	var options Options
	for _, o := range opts {
		o(&options)
	}
	cache := cache.NewRedis(context.Background(), &cache.RedisOpts{
		Host:        config.GetString("redis.addr"),
		Password:    config.GetString("redis.password"),
		Database:    config.GetInt("redis.db"),
		MaxIdle:     10,
		MaxActive:   50,
		IdleTimeout: config.GetInt("redis.idleTimeout"),
	})
	apiV2Wechat := wechat.NewWechat()
	apiV2Wechat.SetCache(cache)

	w := &Wechat{
		options:  &options,
		wechat:   apiV2Wechat,
		cache:    cache,
		payments: make(map[payments.WechatPaySDK]payments.Payment),
	}

	return w
}
