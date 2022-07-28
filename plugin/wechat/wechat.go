package wechat

import (
	"github.com/leor-w/kid/config"

	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/miniprogram"
	miniConf "github.com/silenceper/wechat/v2/miniprogram/config"
	"github.com/silenceper/wechat/v2/pay"
	payConf "github.com/silenceper/wechat/v2/pay/config"
)

type Wechat struct {
	options     *Options
	wechat      *wechat.Wechat
	miniProgram *miniprogram.MiniProgram
	pay         *pay.Pay
	cache       cache.Cache
}

func (w *Wechat) Provide() interface{} {
	if !config.Exist("wechat") {
		panic("not found [wechat] in config")
	}
	return New(
		WithAppid(config.GetString("wechat.appid")),
		WithMiniSecret(config.GetString("wechat.mini.secret")),
		WithPayMchid(config.GetString("wechat.pay.mchid")),
		WithPayKey(config.GetString("wechat.pay.key")),
		WithPayNotifyUrl(config.GetString("wechat.pay.notifyUrl")),
	)
}

type Option func(*Options)

func (w *Wechat) SetCache(cache cache.Cache) {
	w.cache = cache
	w.wechat.SetCache(cache)
}

func (w *Wechat) MiniProgram() *miniprogram.MiniProgram {
	if w.miniProgram == nil {
		w.miniProgram = w.wechat.GetMiniProgram(&miniConf.Config{
			AppID:     w.options.appid,
			AppSecret: w.options.secret,
			Cache:     w.cache,
		})
	}
	return w.miniProgram
}

func (w *Wechat) GetPay() *pay.Pay {
	if w.pay == nil {
		w.pay = w.wechat.GetPay(&payConf.Config{
			AppID:     w.options.appid,
			MchID:     w.options.mchid,
			Key:       w.options.key,
			NotifyURL: w.options.notifyUrl,
		})
	}
	return w.pay
}

func New(opts ...Option) *Wechat {
	var options Options
	for _, o := range opts {
		o(&options)
	}
	c := cache.NewMemory()
	w := wechat.NewWechat()
	w.SetCache(c)
	return &Wechat{
		options: &options,
		wechat:  w,
		cache:   c,
	}
}
