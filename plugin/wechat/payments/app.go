package payments

import (
	"context"

	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"

	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/app"
)

type APP struct {
	service *app.AppApiService
	options *Options
}

func (payment *APP) CloseOrder(ctx context.Context, req *CloseOrderRequest) (*core.APIResult, error) {
	return payment.service.CloseOrder(ctx, app.CloseOrderRequest{
		OutTradeNo: core.String(req.OutTradeNo),
		Mchid:      core.String(payment.options.Mchid),
	})
}

func (payment *APP) Prepay(ctx context.Context, req interface{}) (interface{}, *core.APIResult, error) {
	conf, ok := req.(*app.PrepayRequest)
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

func (payment *APP) QueryOrderById(ctx context.Context, req *QueryOrderByIdRequest) (*payments.Transaction, *core.APIResult, error) {
	return payment.service.QueryOrderById(ctx, app.QueryOrderByIdRequest{
		TransactionId: core.String(req.TransactionId),
		Mchid:         core.String(payment.options.Mchid),
	})
}

func (payment *APP) QueryOrderByTradeNo(ctx context.Context, req *QueryOrderByOutTradeNoRequest) (*payments.Transaction, *core.APIResult, error) {
	return payment.service.QueryOrderByOutTradeNo(ctx, app.QueryOrderByOutTradeNoRequest{
		OutTradeNo: core.String(req.OutTradeNo),
		Mchid:      core.String(payment.options.Mchid),
	})
}

func NewAPP(client *core.Client, opts *Options) *APP {
	return &APP{
		service: &app.AppApiService{Client: client},
		options: opts,
	}
}
