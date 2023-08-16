package payments

import (
	"context"

	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/native"
)

type Native struct {
	service *native.NativeApiService
	options *Options
}

func (payment *Native) CloseOrder(ctx context.Context, req *CloseOrderRequest) (*core.APIResult, error) {
	return payment.service.CloseOrder(ctx, native.CloseOrderRequest{
		OutTradeNo: core.String(req.OutTradeNo),
		Mchid:      core.String(payment.options.Mchid),
	})
}

func (payment *Native) Prepay(ctx context.Context, req interface{}) (interface{}, *core.APIResult, error) {
	conf, ok := req.(*native.PrepayRequest)
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
		Code: *resp.CodeUrl,
	}, result, nil
}

func (payment *Native) QueryOrderById(ctx context.Context, req *QueryOrderByIdRequest) (*payments.Transaction, *core.APIResult, error) {
	return payment.service.QueryOrderById(ctx, native.QueryOrderByIdRequest{
		TransactionId: core.String(req.TransactionId),
		Mchid:         core.String(payment.options.Mchid),
	})
}

func (payment *Native) QueryOrderByTradeNo(ctx context.Context, req *QueryOrderByOutTradeNoRequest) (*payments.Transaction, *core.APIResult, error) {
	return payment.service.QueryOrderByOutTradeNo(ctx, native.QueryOrderByOutTradeNoRequest{
		OutTradeNo: core.String(req.OutTradeNo),
		Mchid:      core.String(payment.options.Mchid),
	})
}

func NewNative(client *core.Client, opts *Options) *Native {
	return &Native{
		service: &native.NativeApiService{Client: client},
		options: opts,
	}
}
