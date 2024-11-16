package pay

type PurchaseVerifyConfig struct {
	/* GooglePlay 包名 */
	PackageName string `json:"packageName"`
	/* GooglePlay 商品 ID */
	ProductID string `json:"productId"`
	/* 内购商品购买的 token */
	PurchaseToken string `json:"purchaseToken"`
}

/*
ProductPurchase 资源指明了用户的应用内商品购买交易的状态。
有关详细信息，请参阅 androidpublisher.purchases.products API 文档。
https://developers.google.com/android-publisher/api-ref/rest/v3/purchases.products?hl=zh-cn#ProductPurchase
*/
type ProductPurchase struct {
	/* 此种类表示 androidpublisher 服务中的 inappPurchase 对象。 */
	Kind string `json:"kind"`
	/* 购买商品的时间，以从公元纪年（1970 年 1 月 1 日）开始计算的毫秒数表示。 */
	PurchaseTimeMillis int64 `json:"purchaseTimeMillis"`
	/* 订单的购买状态。可能的值包括：0。已购买 1.已取消 2.待处理 */
	PurchaseState int64 `json:"purchaseState"`
	/* 应用内商品的使用状态。可能的值包括：0。还将被消耗 1.已使用 */
	ConsumptionState int64 `json:"consumptionState"`
	/* 开发者指定的字符串，其中包含有关订单的补充信息。*/
	DeveloperPayload string `json:"developerPayload"`
	/* 与应用内商品的购买相关联的订单 ID。 */
	OrderId string `json:"orderId"`
	/* 应用内商品的购买类型。仅当此次购买并非通过标准应用内购买结算流程进行时，才会设置此字段。
	可能的值包括：0。测试（即通过许可测试账号购买）1.促销（即使用促销代码购买）。不包括 Play Points 购买交易。2. 激励广告（即通过观看视频广告而不是付费观看） */
	PurchaseType int64 `json:"purchaseType"`
	/* 应用内商品的确认状态。可能的值包括：0。尚未确认 1.已确认 */
	AcknowledgementState int64 `json:"acknowledgementState"`
	/* 系统为标识此次购买交易而生成的购买令牌。可能不存在。 */
	PurchaseToken string `json:"purchaseToken"`
	/* 应用内商品 SKU。可能不存在。*/
	ProductId string `json:"productId"`
	/* 与应用内商品的购买相关联的数量。如果不存在，则数量为 1。 */
	Quantity int64 `json:"quantity"`
	/* 混淆版本的 ID，与您应用中的用户账号唯一关联。
	仅当在购买时使用 https://developer.android.com/reference/com/android/billingclient/api/BillingFlowParams.Builder#setobfuscatedaccountid
	指定时才存在。 */
	ObfuscatedExternalAccountId string `json:"obfuscatedExternalAccountId"`
	/* ID 的混淆版本，与应用中的用户个人资料唯一关联。仅当在购买时使用 https://developer.android.com/reference/com/android/billingclient/api/BillingFlowParams.Builder#setobfuscatedprofileid 指定时才存在。 */
	ObfuscatedExternalProfileId string `json:"obfuscatedExternalProfileId"`
	/* 用户授予产品时的 ISO 3166-1 alpha-2 结算地区代码。 */
	RegionCode string `json:"regionCode"`
	/* 符合退款条件的数量，即尚未退款的数量。该值反映的是基于数量的部分退款和全额退款。 */
	RefundableQuantity int64 `json:"refundableQuantity"`
}
