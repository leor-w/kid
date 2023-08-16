package payments

import (
	"context"

	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
)

type Bridge struct {
	bridge  jsapi.JsapiApiService
	options *Options
}

func (b *Bridge) CloseOrder(ctx context.Context, req *CloseOrderRequest) (*core.APIResult, error) {
	return nil, nil
}

func (b *Bridge) Prepay(ctx context.Context, req interface{}) (interface{}, *core.APIResult, error) {
	conf, ok := req.(*jsapi.PrepayRequest)
	if !ok {
		return nil, nil, ErrorPrepayConf
	}
	if conf.NotifyUrl == nil {
		conf.NotifyUrl = core.String(b.options.PayNotifyUrl)
	}
	if conf.Mchid == nil {
		conf.Mchid = core.String(b.options.Mchid)
	}
	if conf.Appid == nil {
		conf.Appid = core.String(b.options.Appid)
	}
	return b.bridge.PrepayWithRequestPayment(ctx, *conf)
}

func (b *Bridge) QueryOrderById(ctx context.Context,
	req *QueryOrderByIdRequest) (*payments.Transaction, *core.APIResult, error) {
	return nil, nil, nil
}

func (b *Bridge) QueryOrderByTradeNo(ctx context.Context,
	req *QueryOrderByOutTradeNoRequest) (*payments.Transaction, *core.APIResult, error) {
	return nil, nil, nil
}

func NewBridge(client *core.Client, opts *Options) *Bridge {
	return &Bridge{
		bridge: jsapi.JsapiApiService{
			Client: client,
		},
		options: opts,
	}
}
