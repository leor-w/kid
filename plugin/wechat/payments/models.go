package payments

import (
	"net/http"
	"time"
)

// APIResult 微信支付API v3 请求结果
type APIResult struct {
	// 本次请求所使用的 HTTPRequest
	Request *http.Request
	// 本次请求所获得的 HTTPResponse
	Response *http.Response
}

type CloseOrderRequest struct {
	OutTradeNo string
}

type PrepayRequest struct {
	Description      *string    // 商品描述
	OutTradeNo       *string    // 商户订单号
	TimeExpire       *time.Time // 订单失效时间
	Attach           *string    // 附加数据
	GoodsTag         *string    // 回调 URL 地址
	LimitPay         []string   // 商品标记
	SupportFapiao    *bool      // 传入true时，支付成功消息和支付详情页将出现开票入口。需要在微信支付商户平台或微信公众平台开通电子发票功能，传此字段才可生效。
	Total            *int64     // 订单总金额
	Currency         *string    // 支付币种
	CostPrice        *int64     // 1.商户侧一张小票订单可能被分多次支付，订单原价用于记录整张小票的交易金额。 2.当订单原价与支付金额不相等，则不享受优惠。 3.该字段主要用于防止同一张小票分多次支付，以享受多次优惠的情况，正常支付订单不必上传此参数。
	InvoiceId        *string    // 商家小票ID。
	GoodsName        *string    // 商品的实际名称。
	MerchantGoodsId  *string    // 由半角的大小写字母、数字、中划线、下划线中的一种或几种组成。
	Quantity         *int64     // 用户购买的数量。
	UnitPrice        *int64     // 商品单价，单位为分。
	WechatpayGoodsId *string    // 微信支付定义的统一商品编号（没有可不传）。
	ProfitSharing    *bool      // 是否指定分账
	DeviceId         *string    // 商户端设备号
	ClientIp         *string    // 用户终端IP
	StoreAddress     *string    // 详细的商户门店地址
	StoreAreaCode    *string    // 地区编码，详细请见微信支付提供的文档
	StoreId          *string    // 商户侧门店编号
	StoreName        *string    // 商户侧门店名称
	NotifyUrl        *string    // 支付回调地址
}

type QueryOrderByIdRequest struct {
	TransactionId string
}

type QueryOrderByOutTradeNoRequest struct {
	OutTradeNo string
}

type Response struct {
	Code string
}

type CombineTransaction struct {
	Mchid           *string            `json:"mchid,omitempty"`
	TradeType       *string            `json:"trade_type,omitempty"`
	TradeState      *string            `json:"trade_state,omitempty"`
	BankType        *string            `json:"bank_type,omitempty"`
	Attach          *string            `json:"attach,omitempty"`
	SuccessTime     *string            `json:"success_time,omitempty"`
	TransactionId   *string            `json:"transaction_id,omitempty"`
	OutTradeNo      *string            `json:"out_trade_no,omitempty"`
	PromotionDetail []PromotionDetail  `json:"promotion_detail,omitempty"` // 优惠
	Amount          *TransactionAmount `json:"amount,omitempty"`
}

// TransactionAmount
type TransactionAmount struct {
	Currency      *string `json:"currency,omitempty"`
	PayerCurrency *string `json:"payer_currency,omitempty"`
	PayerTotal    *int64  `json:"payer_total,omitempty"`
	Total         *int64  `json:"total,omitempty"`
}

// TransactionPayer
type TransactionPayer struct {
	Openid *string `json:"openid,omitempty"`
}

// PromotionDetail
type PromotionDetail struct {
	// 券ID
	CouponId *string `json:"coupon_id,omitempty"`
	// 优惠名称
	Name *string `json:"name,omitempty"`
	// GLOBAL：全场代金券；SINGLE：单品优惠
	Scope *string `json:"scope,omitempty"`
	// CASH：充值；NOCASH：预充值。
	Type *string `json:"type,omitempty"`
	// 优惠券面额
	Amount *int64 `json:"amount,omitempty"`
	// 活动ID，批次ID
	StockId *string `json:"stock_id,omitempty"`
	// 单位为分
	WechatpayContribute *int64 `json:"wechatpay_contribute,omitempty"`
	// 单位为分
	MerchantContribute *int64 `json:"merchant_contribute,omitempty"`
	// 单位为分
	OtherContribute *int64 `json:"other_contribute,omitempty"`
	// CNY：人民币，境内商户号仅支持人民币。
	Currency    *string                `json:"currency,omitempty"`
	GoodsDetail []PromotionGoodsDetail `json:"goods_detail,omitempty"`
}

// PromotionGoodsDetail
type PromotionGoodsDetail struct {
	// 商品编码
	GoodsId *string `json:"goods_id"`
	// 商品数量
	Quantity *int64 `json:"quantity"`
	// 商品价格
	UnitPrice *int64 `json:"unit_price"`
	// 商品优惠金额
	DiscountAmount *int64 `json:"discount_amount"`
	// 商品备注
	GoodsRemark *string `json:"goods_remark,omitempty"`
}
