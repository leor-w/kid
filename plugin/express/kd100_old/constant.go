package kdn

// ShipperCode 快递公司编码，支持常用快递公司，具体参考快递鸟官网 https://www.yuque.com/kdnjishuzhichi/dfcrg1/mza2ln
type ShipperCode string

const (
	HTKY ShipperCode = "HTKY" // 百世快递
	ZTO              = "ZTO"  // 中通
	STO              = "STO"  // 申通
	YTO              = "YTO"  // 圆通
	YD               = "YD"   // 韵达
	YZPY             = "YZPY" // 邮政快递包裹
	EMS              = "EMS"  // EMS
	SF               = "SF"   // 顺丰
	JD               = "JD"   // 京东快递
	UC               = "UC"   // 优速快递
	DBL              = "DBL"  // 德邦快递
	JTSD             = "JTSD" // 极兔快递
	ZYE              = "ZYE"  // 众邮快递
	ZJS              = "ZJS"  // 宅急送
	FWX              = "FWX"  // 丰网速运
)

// shipperCodeMap 快递公司编码映射
var shipperCodeMap = map[ShipperCode]string{
	HTKY: "百世快递",
	ZTO:  "中通",
	STO:  "申通",
	YTO:  "圆通",
	YD:   "韵达",
	YZPY: "邮政快递包",
	EMS:  "EMS",
	SF:   "顺丰",
	JD:   "京东快递",
	UC:   "优速快递",
	DBL:  "德邦快递",
	JTSD: "极兔快递",
	ZYE:  "众邮快递",
	ZJS:  "宅急送",
	FWX:  "丰网速运",
}

// String 返回快递公司名称
func (sc ShipperCode) String() string {
	return shipperCodeMap[sc]
}

// DetailState 快递细分状态 https://www.yuque.com/kdnjishuzhichi/weubcu/cti4czvp8hnatgue
type DetailState int

const (
	NoTrackInfo               DetailState = 0   // 暂无轨迹信息
	Picked                                = 1   // 已揽收
	OnTheWay                              = 2   // 在途中
	ArriveCity                            = 201 // 到达派件城市
	ArriveTransitCenter                   = 204 // 到达转运中心
	ArrivePoint                           = 205 // 到达派件网点
	SendingPoint                          = 206 // 寄件网点发件
	Dispatching                           = 202 // 派件中
	PlacedCabinet                         = 211 // 已放入快递柜或驿站
	Received                              = 3   // 已签收
	NormalReceived                        = 301 // 正常签收
	AbnormalReceived                      = 302 // 派件异常后最终签收
	BehalfReceived                        = 304 // 代收签收
	ParcelLockerReceived                  = 311 // 快递柜或驿站签收
	Abnormal                              = 4   // 问题件
	NoShippingInfo                        = 401 // 发货无信息
	ExpirationNotSigned                   = 402 // 超时未签收
	UpdatedExpiration                     = 403 // 超时未更新
	Rejected                              = 404 // 拒收(退件)
	DeliveryAnomaly                       = 405 // 派件异常
	ReturnReceipt                         = 406 // 退货签收
	ReturnNotReceived                     = 407 // 退货未签收
	PickupOverdue                         = 412 // 快递柜或驿站超时未取
	Intercepted                           = 413 // 单号已拦截
	Damaged                               = 414 // 破损
	CustomerCancelled                     = 415 // 客户取消发货
	UnableContact                         = 416 // 无法联系
	DeliveryDelay                         = 417 // 配送延迟
	TakenOut                              = 418 // 快件取出
	Redeliver                             = 419 // 重新派送
	ReceivedAddrNotDetailed               = 420 // 收货地址不详细
	RecipientPhoneIncorrect               = 421 // 收件人电话错误
	MisroutedParcel                       = 422 // 错分件
	OutZone                               = 423 // 超区件
	Forwarded                             = 5   // 转寄
	CustomsClearance                      = 6   // 清关
	WaitCustomsClearance                  = 601 // 待清关
	InCustomsClearance                    = 602 // 清关中
	ClearedCustoms                        = 603 // 已清关
	CustomsClearanceException             = 604 // 清关异常
	AwaitingPickup                        = 10  // 待揽件
)

// ExStateMap 快递状态码对应的描述
var ExStateMap = map[DetailState]string{
	NoTrackInfo:               " 暂无轨迹信息",
	Picked:                    " 已揽收",
	OnTheWay:                  " 在途中",
	ArriveCity:                " 到达派件城市",
	ArriveTransitCenter:       " 到达转运中心",
	ArrivePoint:               " 到达派件网点",
	SendingPoint:              " 寄件网点发件",
	Dispatching:               " 派件中",
	PlacedCabinet:             " 已放入快递柜或驿站",
	Received:                  " 已签收",
	NormalReceived:            " 正常签收",
	AbnormalReceived:          " 派件异常后最终签收",
	BehalfReceived:            " 代收签收",
	ParcelLockerReceived:      " 快递柜或驿站签收",
	Abnormal:                  " 问题件",
	NoShippingInfo:            " 发货无信息",
	ExpirationNotSigned:       " 超时未签收",
	UpdatedExpiration:         " 超时未更新",
	Rejected:                  " 拒收(退件)",
	DeliveryAnomaly:           " 派件异常",
	ReturnReceipt:             " 退货签收",
	ReturnNotReceived:         " 退货未签收",
	PickupOverdue:             " 快递柜或驿站超时未取",
	Intercepted:               " 单号已拦截",
	Damaged:                   " 破损",
	CustomerCancelled:         " 客户取消发货",
	UnableContact:             " 无法联系",
	DeliveryDelay:             " 配送延迟",
	TakenOut:                  " 快件取出",
	Redeliver:                 " 重新派送",
	ReceivedAddrNotDetailed:   " 收货地址不详细",
	RecipientPhoneIncorrect:   " 收件人电话错误",
	MisroutedParcel:           " 错分件",
	OutZone:                   " 超区件",
	Forwarded:                 " 转寄",
	CustomsClearance:          " 清关",
	WaitCustomsClearance:      " 待清关",
	InCustomsClearance:        " 清关中",
	ClearedCustoms:            " 已清关",
	CustomsClearanceException: " 清关异常",
	AwaitingPickup:            " 待揽件",
}

// String 返回状态的字符串
func (es DetailState) String() string {
	return ExStateMap[es]
}
