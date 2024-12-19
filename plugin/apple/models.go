package apple

// ApiRetryTimes API 接口请求重试次数
const ApiRetryTimes = 5

const (
	// ProductionBaseURL 苹果支付生产环境基础 URL
	ProductionBaseURL = "https://api.storekit.itunes.apple.com/inApps/v1"
	// SandboxBaseURL 沙箱环境基础 URL
	SandboxBaseURL = "https://api.storekit-sandbox.itunes.apple.com/inApps/v1"
)

// TransactionHistoryQuery 交易历史查询参数
// https://developer.apple.com/documentation/appstoreserverapi/get_transaction_history
type TransactionHistoryQuery struct {
	/**
	 * Revision
	 * 您在查询中使用的令牌，用于为客户请求下一组交易。使用上一次查询返回的值。
	 */
	Revision *string
	/**
	 * StartDate
	 * 可选值，您要查询的开始日期。
	 */
	StartDate *int64
	/**
	 * EndDate
	 * 可选值，您要查询的结束日期。
	 */
	EndDate *int64
	/**
	 * ProductId
	 * 可选值，您要查询的产品的标识符。
	 */
	ProductId *string
	/**
	 * ProductType
	 * 可选值，您要查询的产品类型。
	 */
	ProductType *string
	/**
	 * InAppOwnershipType
	 * 可选值，您要查询的购买类型。
	 */
	InAppOwnershipType *string
	/**
	 * Sort
	 * 可选值，您要查询的排序方式。默认值：ASCENDING。可选值：
	 * - ASCENDING：升序
	 * - DESCENDING：降序
	 */
	Sort *string
	/**
	 * Revoked
	 * 可选值，您要查询的已撤销交易。
	 */
	Revoked *bool
	/**
	 * SubscriptionGroupIdentifier
	 * 可选值，您要查询的订阅组标识符。
	 */
	SubscriptionGroupIdentifier *string
}

// NotificationHistoryQuery 通知历史查询参数
// https://developer.apple.com/documentation/appstoreserverapi/get_notification_history
type NotificationHistoryQuery struct {
	/**
	 * StartDate
	 * 必选值，您要查询的开始日期。需要在当前日期之前的180天之内。
	 */
	StartDate *int64
	/**
	 * EndDate
	 * 必选值，您要查询的结束日期。
	 */
	EndDate *int64
	/**
	 * NotificationType
	 * 可选值，您要查询的通知类型。
	 */
	NotificationType *string
	/**
	 * NotificationSubtype
	 * 可选值，您要查询的通知子类型。
	 */
	NotificationSubtype *string
	/**
	 * OnlyFailures
	 * 可选值，您要查询的失败通知。
	 */
	OnlyFailures *bool
	/**
	 * TransactionId
	 * 可选值，您要查询的交易标识符。
	 */
	TransactionId *string
}

// RefundHistoryQuery 退款历史查询参数
// https://developer.apple.com/documentation/appstoreserverapi/get_refund_history
type RefundHistoryQuery struct {
}

// ConsumptionParams 消费参数
// https://developer.apple.com/documentation/appstoreserverapi/send_consumption_information
type ConsumptionParams struct {
	/**
	 * AccountTenure
	 * 必选值，客户账户的使用年限。
	 */
	AccountTenure *int64
	/**
	 * AppAccountToken
	 * 必选值，完成应用内购买交易的应用内用户帐户的 UUID。
	 */
	AppAccountToken *string
	/**
	 * ConsumptionStatus
	 * 必选值，该值表示客户消费应用内购买的程度。
	 */
	ConsumptionStatus *string
	/**
	 * CustomerConsented
	 * 必选值，表示客户是否同意提供消费数据。
	 */
	CustomerConsented *bool
	/**
	 * DeliveryStatus
	 * 必选值，该值指示应用是否成功交付了正常运行的应用内购买。
	 */
	DeliveryStatus *string
	/**
	 * LifetimeDollarsPurchased
	 * 必选值，该值表示客户在所有平台上在您的应用中进行的应用内购买的总金额（以美元为单位）。
	 */
	LifetimeDollarsPurchased *int64
	/**
	 * LifetimeDollarsRefunded
	 * 必选值，该值表示客户在您的应用中在所有平台上收到的退款总金额（以美元为单位）。
	 */
	LifetimeDollarsRefunded *int64
	/**
	 * Platform
	 * 必选值，该值指示客户进行应用内购买的平台。
	 */
	Platform *string
	/**
	 * PlayTime
	 * 必选值，表示客户使用该应用程序的时间量的值。
	 */
	PlayTime *int64
	/**
	 * RefundPreference
	 * 可选值，该值表示根据您的操作逻辑，Apple 是否应该批准退款的偏好。
	 */
	RefundPreference *bool
	/**
	 * SampleContentProvided
	 * 必选值，表示您是否在购买之前提供了内容的免费样本或试用版，或者有关其功能的信息。
	 */
	SampleContentProvided *bool
	/**
	 * UserStatus
	 * 必选值，客户账户的状态。
	 */
	UserStatus *string
}
