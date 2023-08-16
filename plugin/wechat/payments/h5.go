package payments

import (
	"context"

	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/h5"
)

type H5Payment struct {
	service *h5.H5ApiService
	options *Options
}

func (payment *H5Payment) CloseOrder(ctx context.Context, req *CloseOrderRequest) (*core.APIResult, error) {
	return payment.service.CloseOrder(ctx, h5.CloseOrderRequest{
		OutTradeNo: core.String(req.OutTradeNo),
		Mchid:      core.String(payment.options.Mchid),
	})
}

func (payment *H5Payment) Prepay(ctx context.Context, req interface{}) (interface{}, *core.APIResult, error) {
	conf, ok := req.(*h5.PrepayRequest)
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
		Code: *resp.H5Url,
	}, result, nil
}

func (payment *H5Payment) QueryOrderById(ctx context.Context, req *QueryOrderByIdRequest) (*payments.Transaction, *core.APIResult, error) {
	return payment.service.QueryOrderById(ctx, h5.QueryOrderByIdRequest{
		TransactionId: core.String(req.TransactionId),
		Mchid:         core.String(payment.options.Mchid),
	})
}

func (payment *H5Payment) QueryOrderByTradeNo(ctx context.Context, req *QueryOrderByOutTradeNoRequest) (*payments.Transaction, *core.APIResult, error) {
	return payment.service.QueryOrderByOutTradeNo(ctx, h5.QueryOrderByOutTradeNoRequest{
		OutTradeNo: core.String(req.OutTradeNo),
		Mchid:      core.String(payment.options.Mchid),
	})
}

func NewH5(client *core.Client, opts *Options) *H5Payment {
	return &H5Payment{
		service: &h5.H5ApiService{Client: client},
		options: opts,
	}
}
