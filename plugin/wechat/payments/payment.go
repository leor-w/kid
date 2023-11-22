package payments

import (
	"context"
	"errors"

	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"

	"github.com/wechatpay-apiv3/wechatpay-go/core"
)

type WechatPaySDK string

const (
	PaymentTypeAPP    WechatPaySDK = "APP"    // app 支付
	PaymentTypeJSAPI  WechatPaySDK = "JSAPI"  // jsapi 支付 或 小程序支付
	PaymentTypeH5     WechatPaySDK = "H5"     // h5 支付
	PaymentTypeNative WechatPaySDK = "NATIVE" // native 网页支付
)

var (
	NotSupportPayApiType = errors.New("未支持的支付类型")
	ErrorPrepayConf      = errors.New("错误的配置类型")
)

type Payment interface {
	CloseOrder(context.Context, *CloseOrderRequest) (*core.APIResult, error)
	Prepay(context.Context, interface{}) (interface{}, *core.APIResult, error)
	QueryOrderById(context.Context, *QueryOrderByIdRequest) (*payments.Transaction, *core.APIResult, error)
	QueryOrderByTradeNo(context.Context, *QueryOrderByOutTradeNoRequest) (*payments.Transaction, *core.APIResult, error)
}

type Option func(*Options)

func New(apiType WechatPaySDK, client *core.Client, opts ...Option) (Payment, error) {
	var options Options
	for i := range opts {
		opts[i](&options)
	}
	switch apiType {
	case PaymentTypeAPP:
		return NewAPP(client, &options), nil
	case PaymentTypeJSAPI:
		return NewBridge(client, &options), nil
	case PaymentTypeH5:
		return NewH5(client, &options), nil
	case PaymentTypeNative:
		return NewNative(client, &options), nil
	default:
		return nil, NotSupportPayApiType
	}
}
