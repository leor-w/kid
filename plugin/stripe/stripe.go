package stripe

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/lsjhtang/kid/config"
	"github.com/lsjhtang/kid/plugin"
	"github.com/lsjhtang/kid/utils"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/paymentintent"
	"github.com/stripe/stripe-go/v81/paymentlink"
	"github.com/stripe/stripe-go/v81/price"
	"github.com/stripe/stripe-go/v81/webhook"
)

type Stripe struct {
	options *Options
}

type Option func(*Options)

func (s *Stripe) Provide(ctx context.Context) any {
	var configName string
	name, ok := ctx.Value(new(plugin.NameKey)).(string)
	if ok && len(name) > 0 {
		configName = "." + name
	}
	confPrefix := fmt.Sprintf("stripe%s", configName)
	if !config.Exist(confPrefix) {
		panic(fmt.Sprintf("配置文件未找到配置项 [%s]", confPrefix))
	}
	return NewStripe(
		WithSecretKey(config.GetString(utils.GetConfigurationItem(confPrefix, "secret_key"))),
		WithWebhookSecret(config.GetString(utils.GetConfigurationItem(confPrefix, "webhook_secret"))),
		WithRedirectType(config.GetString(utils.GetConfigurationItem(confPrefix, "redirect_type"))),
		WithRedirectDomain(config.GetString(utils.GetConfigurationItem(confPrefix, "redirect_domain"))),
	)
}

func (s *Stripe) BuildPaymentLinkURL(conf *BuildPaymentLinkConfig) (string, error) {
	var afterCompletion *stripe.PaymentLinkAfterCompletionParams
	var redirectUrl string
	if len(conf.RedirectURI) > 0 {
		afterCompletion = &stripe.PaymentLinkAfterCompletionParams{
			Type: stripe.String(s.options.RedirectType),
		}
		redirectUrl = s.options.RedirectDomain + conf.RedirectURI
		if s.options.RedirectType == RedirectTypeRedirect {
			afterCompletion.Redirect = &stripe.PaymentLinkAfterCompletionRedirectParams{
				URL: stripe.String(redirectUrl),
			}
		} else if s.options.RedirectType == RedirectTypeHostedConfirmation {
			afterCompletion.HostedConfirmation = &stripe.PaymentLinkAfterCompletionHostedConfirmationParams{
				CustomMessage: stripe.String(redirectUrl),
			}
		}
	}

	params := &stripe.PaymentLinkParams{
		AfterCompletion: afterCompletion,
		LineItems: []*stripe.PaymentLinkLineItemParams{
			{
				Price:    stripe.String(conf.Price),
				Quantity: stripe.Int64(1),
			},
		},
		Metadata: conf.Metadata,
	}
	pl, err := paymentlink.New(params)
	if err != nil {
		return "", fmt.Errorf("创建支付链接失败: %s", err.Error())
	}
	return pl.URL, nil

}

// GenerateAppPayment 为移动应用创建支付Intent
// 返回clientSecret，可以在移动应用中使用该值初始化支付流程
func (s *Stripe) GenerateAppPayment(conf *AppPaymentConfig) (*stripe.PaymentIntent, error) {

	if conf.Amount <= 0 {
		return nil, errors.New("付款金额必须大于0")
	}

	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(conf.Amount),
		Currency: stripe.String(string(conf.Currency)),
		// 设置自动支付方法，Stripe会自动选择支付方法
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
		Description: stripe.String(conf.Description),
		Metadata:    conf.Metadata,
	}

	// 如果配置了客户ID，添加到参数中
	if conf.Customer != "" {
		params.Customer = stripe.String(conf.Customer)
	}

	// 如果配置了收据邮箱，添加到参数中
	if conf.ReceiptEmail != "" {
		params.ReceiptEmail = stripe.String(conf.ReceiptEmail)
	}

	// 创建PaymentIntent
	paymentIntent, err := paymentintent.New(params)
	if err != nil {
		return nil, fmt.Errorf("创建支付意向失败: %s", err.Error())
	}

	return paymentIntent, nil
}

// VerifyAppPaymentWebhook 验证移动应用支付的Webhook回调
// 当移动应用完成支付后，Stripe会发送一个webhook到服务器
// 该方法用于验证webhook的签名并返回PaymentIntent对象
func (s *Stripe) VerifyAppPaymentWebhook(req *http.Request) (*stripe.PaymentIntent, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, fmt.Errorf("读取请求体失败: %s", err.Error())
	}
	event, err := webhook.ConstructEvent(body, req.Header.Get("Stripe-Signature"), s.options.WebhookSecret)
	if err != nil {
		return nil, fmt.Errorf("验证签名失败: %s", err.Error())
	}
	// 检查事件类型是否为支付成功
	if event.Type != InAppPurchaseSuccess {
		return nil, fmt.Errorf("事件类型错误: %s", event.Type)
	}
	var paymentIntent *stripe.PaymentIntent
	err = json.Unmarshal(event.Data.Raw, &paymentIntent)
	if err != nil {
		return nil, fmt.Errorf("解析事件数据失败: %s", err.Error())
	}
	return paymentIntent, nil
}

// GetPaymentIntent 通过ID获取PaymentIntent对象
// 可用于查询支付状态
func (s *Stripe) GetPaymentIntent(id string) (*stripe.PaymentIntent, error) {
	paymentIntent, err := paymentintent.Get(id, nil)
	if err != nil {
		return nil, fmt.Errorf("获取支付意向失败: %s", err.Error())
	}
	return paymentIntent, nil
}

// GetPrice 通过ID获取Price对象
func (s *Stripe) GetPrice(id string) (*stripe.Price, error) {
	price, err := price.Get(id, nil)
	if err != nil {
		return nil, fmt.Errorf("获取价格失败: %s", err.Error())
	}
	return price, nil
}

func NewStripe(opts ...Option) *Stripe {
	options := &Options{}
	for _, opt := range opts {
		opt(options)
	}
	stripe.Key = options.SecretKey
	return &Stripe{
		options: options,
	}
}
