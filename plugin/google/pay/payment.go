package pay

import (
	"context"
	"fmt"

	"github.com/leor-w/kid/utils"
	"google.golang.org/api/androidpublisher/v3"
	"google.golang.org/api/option"

	"github.com/leor-w/kid/config"
	"github.com/leor-w/kid/plugin"
)

type GooglePay struct {
	options *Options
}

type Option func(*Options)

func (gp *GooglePay) Provide(ctx context.Context) any {
	var configName string
	name, ok := ctx.Value(new(plugin.NameKey)).(string)
	if ok && len(name) > 0 {
		configName = "." + name
	}
	confPrefix := fmt.Sprintf("google.pay%s", configName)
	if !config.Exist(fmt.Sprintf(confPrefix)) {
		panic(fmt.Sprintf("配置文件未找到配置项 [%s]", confPrefix))
	}
	return NewGooglePay(
		WithServerAccountFile(config.GetString(utils.GetConfigurationItem(confPrefix, "server_account_file"))),
	)
}

func (gp *GooglePay) VerifyInAPPPurchase(purchaseConf *PurchaseVerifyConfig) (*ProductPurchase, error) {
	srv, err := androidpublisher.NewService(context.Background(), option.WithCredentialsFile(gp.options.ServerAccountFile))
	if err != nil {
		return nil, fmt.Errorf("google 内购服务初始化失败: %s", err.Error())
	}
	// 验证内购
	purchase, err := srv.Purchases.Products.Get(purchaseConf.PackageName, purchaseConf.ProductID, purchaseConf.PurchaseToken).Do()
	if err != nil {
		return nil, fmt.Errorf("google 内购验证失败: %s", err.Error())
	}
	if purchase.PurchaseState != 0 {
		return nil, fmt.Errorf("google 内购状态异常: %d", purchase.PurchaseState)
	}
	return &ProductPurchase{
		Kind:                        purchase.Kind,
		PurchaseTimeMillis:          purchase.PurchaseTimeMillis,
		PurchaseState:               purchase.PurchaseState,
		ConsumptionState:            purchase.ConsumptionState,
		DeveloperPayload:            purchase.DeveloperPayload,
		OrderId:                     purchase.OrderId,
		PurchaseType:                *purchase.PurchaseType,
		AcknowledgementState:        purchase.AcknowledgementState,
		PurchaseToken:               purchase.PurchaseToken,
		ProductId:                   purchase.ProductId,
		Quantity:                    purchase.Quantity,
		ObfuscatedExternalAccountId: purchase.ObfuscatedExternalAccountId,
		ObfuscatedExternalProfileId: purchase.ObfuscatedExternalProfileId,
		RegionCode:                  purchase.RegionCode,
		RefundableQuantity:          purchase.RefundableQuantity,
	}, nil
}

func NewGooglePay(opts ...Option) *GooglePay {
	options := &Options{}
	for _, o := range opts {
		o(options)
	}
	return &GooglePay{options: options}
}
