package stripe

const (
	InAppPurchaseSuccess = "payment_intent.succeeded"
	InAppPurchaseFail    = "payment_intent.payment_failed"
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

const (
	RedirectTypeRedirect           = "redirect"
	RedirectTypeHostedConfirmation = "hosted_confirmation"
)
