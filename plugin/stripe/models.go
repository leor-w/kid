package stripe

const (
	InAppPurchaseSuccess = "payment_intent.succeeded" // 支付成功
	InAppPurchaseFail    = "payment_intent.payment_failed" // 支付失败
	InAppPurchaseAttached = "payment_method.attached" // 支付方法附加成功
)

const (
	CheckoutSeesionCompleted = "checkout.session.completed"
)

const (
	CheckoutSeesionPaid              = "paid"                // 已支付
	CheckoutSeesionNoPaymentRequired = "no_payment_required" // 无需支付
	CheckoutSeesionUnpaid            = "unpaid"              // 未支付
)

type BuildPaymentLinkConfig struct {
	TradeNo     string
	Price       string
	Quantity    int
	Metadata    map[string]string
	RedirectURI string
}

// AppPaymentConfig 用于移动应用程序支付的配置
type AppPaymentConfig struct {
	Amount      int64             // 支付金额（以分为单位）
	Currency    string            // 货币代码，例如 "usd", "cny"
	Description string            // 订单描述
	Metadata    map[string]string // 元数据，可以包含订单号等信息
	Customer    string            // 可选，客户ID
	ReceiptEmail string           // 可选，收据邮箱
}

const (
	RedirectTypeRedirect           = "redirect"
	RedirectTypeHostedConfirmation = "hosted_confirmation"
)
