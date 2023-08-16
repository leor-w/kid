package payments

import (
	"context"

	"github.com/wechatpay-apiv3/wechatpay-go/core"

	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
)

type Jsapi struct {
	service *jsapi.JsapiApiService
	options *Options
}

func (payment *Jsapi) CloseOrder(ctx context.Context, req *CloseOrderRequest) (*core.APIResult, error) {
	return payment.service.CloseOrder(ctx, jsapi.CloseOrderRequest{
		OutTradeNo: core.String(req.OutTradeNo),
		Mchid:      core.String(payment.options.Mchid),
	})
}

func (payment *Jsapi) Prepay(ctx context.Context, req interface{}) (interface{}, *core.APIResult, error) {
	conf, ok := req.(*jsapi.PrepayRequest)
	if !ok {
		return nil, nil, ErrorPrepayConf
	}
	if conf.NotifyUrl == nil {
		conf.NotifyUrl = core.String(payment.options.PayNotifyUrl)
	}
	if conf.Appid == nil {
		conf.Appid = core.String(payment.options.Appid)
	}
	if conf.Mchid == nil {
		conf.Mchid = core.String(payment.options.Mchid)
	}
	resp, result, err := payment.service.Prepay(ctx, *conf)
	if err != nil {
		return nil, nil, err
	}
	return &Response{
		Code: *resp.PrepayId,
	}, result, nil
}

func (payment *Jsapi) QueryOrderById(ctx context.Context, req *QueryOrderByIdRequest) (*payments.Transaction, *core.APIResult, error) {
	return payment.service.QueryOrderById(ctx, jsapi.QueryOrderByIdRequest{
		TransactionId: core.String(req.TransactionId),
		Mchid:         core.String(payment.options.Mchid),
	})
}

func (payment *Jsapi) QueryOrderByTradeNo(ctx context.Context, req *QueryOrderByOutTradeNoRequest) (*payments.Transaction, *core.APIResult, error) {
	return payment.service.QueryOrderByOutTradeNo(ctx, jsapi.QueryOrderByOutTradeNoRequest{
		OutTradeNo: core.String(req.OutTradeNo),
		Mchid:      core.String(payment.options.Mchid),
	})
}

func NewJsapi(client *core.Client, opts *Options) *Jsapi {
	return &Jsapi{
		service: &jsapi.JsapiApiService{Client: client},
		options: opts,
	}
}
