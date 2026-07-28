package pay

import (
	"context"
	"fmt"

	"google.golang.org/api/androidpublisher/v3"
	"google.golang.org/api/option"

	"github.com/lsjhtang/kid/config"
	"github.com/lsjhtang/kid/plugin"
	"github.com/lsjhtang/kid/utils"
)

type GooglePay struct {
	googleService *androidpublisher.Service
	options       *Options
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
		panic(fmt.Sprintf("google play: 配置文件未找到配置项 [%s]", confPrefix))
	}
	return NewGooglePay(
		WithServerAccountFile(config.GetString(utils.GetConfigurationItem(confPrefix, "server_account_file"))),
	)
}

// productGet 获取内购商品订单信息
func (gp *GooglePay) productGet(conf *PurchaseConfig) (*androidpublisher.ProductPurchase, error) {
	srv, err := androidpublisher.NewService(context.Background(), option.WithCredentialsFile(gp.options.ServerAccountFile))
	if err != nil {
		return nil, fmt.Errorf("google play: 内购服务初始化失败: %s", err.Error())
	}
	// 通过客户端或者Google推送的购买信息验证购买状态
	purchase, err := srv.Purchases.Products.Get(conf.PackageName, conf.ProductID, conf.PurchaseToken).Do()
	if err != nil {
		return nil, fmt.Errorf("google play: 内购商品订单验证失败: %s", err.Error())
	}
	return purchase, nil
}

// subscriptionGet 获取内购订阅商品订单信息
func (gp *GooglePay) subscriptionGet(conf *PurchaseConfig) (*androidpublisher.SubscriptionPurchase, error) {
	srv, err := androidpublisher.NewService(context.Background(), option.WithCredentialsFile(gp.options.ServerAccountFile))
	if err != nil {
		return nil, fmt.Errorf("google play: 内购服务初始化失败: %s", err.Error())
	}
	// 通过客户端或者Google推送的购买信息验证购买状态
	purchase, err := srv.Purchases.Subscriptions.Get(conf.PackageName, conf.ProductID, conf.PurchaseToken).Do()
	if err != nil {
		return nil, fmt.Errorf("google play: 内购商品订单验证失败: %s", err.Error())
	}
	return purchase, nil
}

// GetAndConfirmConsumable 获取并确认消耗型商品，如果内购商品已消耗，无法重复消耗。如果内购商品未消耗，确认消耗
func (gp *GooglePay) GetAndConfirmConsumable(purchaseConf *PurchaseConfig) (*ProductPurchase, error) {
	purchase, err := gp.GetConsumable(purchaseConf)
	if err != nil {
		return nil, err
	}
	if err := gp.ConfirmConsumable(purchaseConf); err != nil {
		return nil, err
	}
	return purchase, nil
}

// GetConsumable 获取内购消耗型商品
func (gp *GooglePay) GetConsumable(purchaseConf *PurchaseConfig) (*ProductPurchase, error) {
	purchase, err := gp.productGet(purchaseConf)
	if err != nil {
		return nil, err
	}
	// 验证商品购买的状态，必须为已购买状态
	if purchase.PurchaseState != 0 {
		return nil, fmt.Errorf("google play: 内购商品支付状态错误，当前状态为: %d", purchase.PurchaseState)
	}
	// 如果内购商品已消耗，无法重复消耗
	if purchaseConf.IsConsumed && purchase.ConsumptionState == 1 {
		return nil, fmt.Errorf("google play: 内购商品已消耗，无法重复消耗")
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

// ConfirmConsumable 确认消耗型商品
func (gp *GooglePay) ConfirmConsumable(conf *PurchaseConfig) error {
	srv, err := androidpublisher.NewService(context.Background(), option.WithCredentialsFile(gp.options.ServerAccountFile))
	if err != nil {
		return fmt.Errorf("google play: 内购服务初始化失败: %s", err.Error())
	}
	// 向 google 确认已消耗内购，如果3天内未消耗，google 会自动退款
	if err := srv.Purchases.Products.Consume(conf.PackageName, conf.ProductID, conf.PurchaseToken).Do(); err != nil {
		return fmt.Errorf("google play: 内购商品确认消耗失败: %s", err.Error())
	}
	return nil
}

// GetNonConsumable 获取非消耗型商品
func (gp *GooglePay) GetNonConsumable(conf *PurchaseConfig) (*ProductPurchase, error) {
	// 通过客户端或者Google推送的购买信息验证购买状态
	purchase, err := gp.googleService.Purchases.Products.Get(conf.PackageName, conf.ProductID, conf.PurchaseToken).Do()
	if err != nil {
		return nil, fmt.Errorf("google play: 内购验证失败: %s", err.Error())
	}
	// 验证商品购买的状态，必须为已购买状态
	if purchase.PurchaseState != 0 {
		return nil, fmt.Errorf("google play: 内购商品支付状态错误，当前状态为: %d", purchase.PurchaseState)
	}
	// 如果内购商品已确认，无法重复确认
	if purchase.AcknowledgementState == 1 {
		return nil, fmt.Errorf("google play: 内购商品已确认，无法重复确认")
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

func (gp *GooglePay) ConfirmNonConsumable(conf *PurchaseConfig) error {
	// 确认非消耗型商品
	if err := gp.googleService.Purchases.Products.
		Acknowledge(
			conf.PackageName,
			conf.ProductID,
			conf.PurchaseToken,
			nil,
		).Do(); err != nil {
		return fmt.Errorf("google play: 内购商品确认失败: %s", err.Error())
	}
	return nil
}

// VerifySubscriptionsPurchase 确认内购订阅商品
func (gp *GooglePay) VerifySubscriptionsPurchase(conf *PurchaseConfig) error {
	// 获取内购订阅商品订单信息
	purchase, err := gp.googleService.Purchases.Subscriptionsv2.Get(conf.PackageName, conf.PurchaseToken).Do()
	if err != nil {
		return fmt.Errorf("google play: 内购订阅商品验证失败: %s", err.Error())
	}
	// 验证商品购买的状态，必须为已购买状态
	if purchase.SubscriptionState != SubscriptionStateActive {
		return fmt.Errorf("google play: 内购订阅商品支付状态错误，当前状态为: %v", purchase.SubscriptionState)
	}
	// 如果内购订阅商品已确认，无法重复确认
	if purchase.AcknowledgementState == AcknowledgementStateAcknowledged {
		return fmt.Errorf("google play: 内购订阅商品已确认，无法重复确认")
	}
	// 确认新的订阅商品
	if err := gp.googleService.Purchases.Subscriptions.
		Acknowledge(conf.PackageName, conf.ProductID, conf.PurchaseToken, nil).
		Do(); err != nil {
		return fmt.Errorf("google play: 内购订阅商品确认失败: %s", err.Error())
	}
	return nil
}

func NewGooglePay(opts ...Option) *GooglePay {
	options := &Options{}
	for _, o := range opts {
		o(options)
	}
	googleService, err := androidpublisher.NewService(context.Background(), option.WithCredentialsFile(options.ServerAccountFile))
	if err != nil {
		panic(fmt.Sprintf("google play: 内购服务初始化失败: %s", err.Error()))
	}
	return &GooglePay{
		googleService: googleService,
		options:       options,
	}
}
