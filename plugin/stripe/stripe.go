package stripe

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/leor-w/kid/config"
	"github.com/leor-w/kid/plugin"
	"github.com/leor-w/kid/utils"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/paymentlink"
	"github.com/stripe/stripe-go/v81/webhook"
	"io"
	"net/http"
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
	if !config.Exist(fmt.Sprintf(confPrefix)) {
		panic(fmt.Sprintf("配置文件未找到配置项 [%s]", confPrefix))
	}
	return NewStripe(
		WithSecretKey(config.GetString(utils.GetConfigurationItem(confPrefix, "secret_key"))),
		WithWebhookSecret(config.GetString(utils.GetConfigurationItem(confPrefix, "webhook_secret"))),
	)
}

func (s *Stripe) BuildPaymentLinkURL(conf *BuildPaymentLinkConfig) (string, error) {
	params := &stripe.PaymentLinkParams{
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

func (s *Stripe) VerifyPaymentIntentWebhook(req *http.Request) (*stripe.CheckoutSession, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, fmt.Errorf("读取请求体失败: %s", err.Error())
	}
	event, err := webhook.ConstructEvent(body, req.Header.Get("Stripe-Signature"), s.options.WebhookSecret)
	if err != nil {
		return nil, fmt.Errorf("验证签名失败: %s", err.Error())
	}
	if event.Type != CheckoutSeesionCompleted {
		return nil, fmt.Errorf("事件类型错误: %s", event.Type)
	}
	var paymentIntent *stripe.CheckoutSession
	err = json.Unmarshal(event.Data.Raw, &paymentIntent)
	if err != nil {
		return nil, fmt.Errorf("解析事件数据失败: %s", err.Error())
	}
	return paymentIntent, nil
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
