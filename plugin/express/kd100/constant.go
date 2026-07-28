package kd100

import "github.com/lsjhtang/kid/plugin/express"

const (
	ResultV2 = "resultv2"
)

const (
	SubscribeStatusPolling   = "polling"   // 监控中
	SubscribeStatusShutdown  = "shutdown"  // 结束
	SubscribeStatusAbort     = "abort"     // 中止 当message为“3天查询无记录”或“60天无变化时”status= abort
	SubscribeStatusUpdateAll = "updateall" // 重新推送
)

type TrackStatus string

const (
	StatusInTransit   TrackStatus = "0"    // 在途
	StatusArrivedCity TrackStatus = "1001" // 到达派件城市
	StatusMainline    TrackStatus = "1002" // 运输干线
	StatusTransfer    TrackStatus = "1003" // 转递

	StatusAwaitPickup    TrackStatus = "1"   // 揽收
	StatusOrderPlaced    TrackStatus = "101" // 已下单
	StatusAwaitingPickup TrackStatus = "102" // 待揽收
	StatusPickedUp       TrackStatus = "103" // 已揽收

	StatusDifficult         TrackStatus = "2"   // 疑难件
	StatusSignExpired       TrackStatus = "201" // 签收超时
	StatusUpdateExpired     TrackStatus = "202" // 更新超时
	StatusRefuseReceive     TrackStatus = "203" // 拒收（退件）
	StatusDeliveryException TrackStatus = "204" // 派件异常
	StatusStationException  TrackStatus = "205" // 站点或快递柜超时未取件
	StatusNotContacted      TrackStatus = "206" // 未联系上收件人
	StatusOutDeliveryArea   TrackStatus = "207" // 超出配送区域
	StatusUndelivered       TrackStatus = "208" // 滞留 – 暂无法进行派件
	StatusDamaged           TrackStatus = "209" // 破损
	StatusCancelOrder       TrackStatus = "210" // 用户取消订单

	StatusSignedDelivery TrackStatus = "3"   // 签收
	StatusSelfSigned     TrackStatus = "301" // 本人签收
	StatusAbnormalSigned TrackStatus = "302" // 派件异常后签收
	StatusOtherSigned    TrackStatus = "303" // 代签
	StatusStationSigned  TrackStatus = "304" // 站点代签或者投柜代签

	StatusReturnSigned TrackStatus = "4"   // 退签
	StatusSalesOrder   TrackStatus = "401" // 快递单已撤销

	StatusDelivering       TrackStatus = "5"   // 派件
	StatusDeliveredStation TrackStatus = "501" // 投柜或驿站 快递已投柜或已送达驿站

	StatusReturned TrackStatus = "6" // 退回

	StatusTransferAnother TrackStatus = "7" // 转寄

	StatusCustomsClearance          TrackStatus = "8"  // 清关
	StatusWaitCustomsClearance      TrackStatus = "10" // 等待清关
	StatusCustomsCleaning           TrackStatus = "11" // 清关中
	StatusCustomsClearanceCompleted TrackStatus = "12" // 清关完成
	StatusCustomsClearanceFail      TrackStatus = "13" // 清关异常

	StatusRefusal TrackStatus = "14" // 拒签
)

var statusToStandardState = map[TrackStatus]express.State{
	StatusAwaitPickup:    express.StatePickUp,
	StatusOrderPlaced:    express.StatePickUp,
	StatusAwaitingPickup: express.StatePickUp,
	StatusPickedUp:       express.StatePickUp,

	StatusInTransit:   express.StateInTransit,
	StatusArrivedCity: express.StateInTransit,
	StatusMainline:    express.StateInTransit,
	StatusTransfer:    express.StateInTransit,

	StatusDelivering:       express.StateDispatch,
	StatusDeliveredStation: express.StateDispatch,

	StatusSignedDelivery: express.StateSign,
	StatusSelfSigned:     express.StateSign,
	StatusAbnormalSigned: express.StateSign,
	StatusOtherSigned:    express.StateSign,
	StatusStationSigned:  express.StateSign,

	StatusDifficult:         express.StateDifficult,
	StatusSignExpired:       express.StateDifficult,
	StatusUpdateExpired:     express.StateDifficult,
	StatusRefuseReceive:     express.StateDifficult,
	StatusDeliveryException: express.StateDifficult,
	StatusStationException:  express.StateDifficult,
	StatusNotContacted:      express.StateDifficult,
	StatusOutDeliveryArea:   express.StateDifficult,
	StatusUndelivered:       express.StateDifficult,
	StatusDamaged:           express.StateDifficult,
	StatusCancelOrder:       express.StateDifficult,

	StatusReturnSigned: express.StateReject,
	StatusSalesOrder:   express.StateReject,

	StatusReturned: express.StateReturned,

	StatusTransferAnother: express.StateTransferAnother,

	StatusCustomsClearance:          express.StateClearance,
	StatusWaitCustomsClearance:      express.StateClearance,
	StatusCustomsCleaning:           express.StateClearance,
	StatusCustomsClearanceCompleted: express.StateClearance,
	StatusCustomsClearanceFail:      express.StateClearance,

	StatusRefusal: express.StateRefuse,
}

func (ts TrackStatus) GetStatusToStandardState() express.State {
	return statusToStandardState[ts]
}

func GetStatusToStandardState(status string) express.State {
	return statusToStandardState[TrackStatus(status)]
}

type Sort string

const (
	SortDesc = "desc" // 降序
	SortAsc  = "asc"  // 升序
)

// MonitorStatus 监控状态
type MonitorStatus uint8

const (
	Polling   MonitorStatus = iota // 监控中
	Shutdown                       // 结束
	Abort                          // 中止
	UpdateAll                      // 重新推送
)

var MonitorStatusMap = map[MonitorStatus]string{
	Polling:   "polling",
	Shutdown:  "shutdown",
	Abort:     "abort",
	UpdateAll: "updateall",
}

var MonitorStatusReverseMap = map[string]MonitorStatus{
	"polling":   Polling,
	"shutdown":  Shutdown,
	"abort":     Abort,
	"updateall": UpdateAll,
}

func (s MonitorStatus) String() string {
	return MonitorStatusMap[s]
}

func ReverseString(monitor string) MonitorStatus {
	return MonitorStatusReverseMap[monitor]
}

var PayTypeMap = map[express.PayType]string{
	express.PayTypeSender:    "SHIPPER",
	express.PayTypeRecipient: "CONSIGNEE",
	express.PayTypeMonthly:   "MONTHLY",
}

func GetPayType(payType express.PayType) string {
	return PayTypeMap[payType]
}

var PayTypeReverseMap = map[string]express.PayType{
	"SHIPPER":   express.PayTypeSender,
	"CONSIGNEE": express.PayTypeRecipient,
	"MONTHLY":   express.PayTypeMonthly,
}

func GetPayTypeReverse(payType string) express.PayType {
	return PayTypeReverseMap[payType]
}

// CardType 证件类型
type CardType uint8

const (
	CardTypeIdCard        CardType = iota + 1 // 身份证
	CardTypeHKMacaoPermit                     // 港澳通行证
	CardTypeTWPermit                          // 台湾通行证
	CardTypePassport                          // 护照
)

// ExpressType 快递类型
type ExpressType uint8

const (
	ExpressTypeDomesiticExpress       = iota + 1 // 国内快递
	ExpressTypeInternationalExpress              // 国际快递
	ExpressTypeDomesiticLogistics                // 国内物流
	ExpressTypeInternationalLogistics            // 国际物流
)

type FeeType uint8

func (f FeeType) String() string {
	return FeeTypeMap[f]
}

func ParseFeeType(feeType string) FeeType {
	return FeeTypeReverseMap[feeType]
}

func (f FeeType) GetFeeTypeDesc() string {
	return FeeTypeDescMap[f]
}

func GetFeeTypeDesc(feeType string) string {
	return FeeTypeReverseDescMap[feeType]
}

const (
	FeeTypeInsuranceFee                FeeType = iota + 1 // 保价费
	FeeTypePackagingFee                                   // 包装费
	FeeTypeCollectionFee                                  // 代收货款
	FeeTypeSignAndReturn                                  // 签单返还
	FeeTypeHandlingCharges                                // 装卸服务
	FeeTypePickUpServiceFee                               // 提货服务
	FeeTypeDeliveryServiceFee                             // 配送服务
	FeeTypeResourceConditioningCharges                    // 资源调节费
	FeeTypeAnimalQuarantineCertificate                    // 动检证
	FeeTypeSpecialWarehousingFee                          // 特殊入仓
	FeeTypeSurcharge                                      // 附加费
	FeeTypeExpressCompensation                            // 快递赔付费
	FeeTypeReceiveDeliveryFeeClaims                       // 收派服务费-理赔
	FeeTypeExpressAddrChangeFee                           // 快递改址费
	FeeTypeExceededAreaFee                                // 超区服务费
	FeeTypeOverLengthOverWeightFee                        // 超长超重附加费
	FeeTypeReturnFee                                      // 逆向费用
	FeeTypeFreshFee                                       // 保鲜服务费
	FeeTypeFullInsuranceFee                               // 顺丰足额保
	FeeTypeOtherFee                                       // 其他费用
)

var FeeTypeMap = map[FeeType]string{
	FeeTypeInsuranceFee:                "INSURANCEFEE",
	FeeTypePackagingFee:                "PACKAGINGFEE",
	FeeTypeCollectionFee:               "COLLECTIONFEE",
	FeeTypeSignAndReturn:               "SIGNANDRETURN",
	FeeTypeHandlingCharges:             "HANDLINGCHARGES",
	FeeTypePickUpServiceFee:            "PICKUPSERVICEFEE",
	FeeTypeDeliveryServiceFee:          "DELIVERYSERVICEFEE",
	FeeTypeResourceConditioningCharges: "RESOURCECONDITIONINGCHARGES",
	FeeTypeAnimalQuarantineCertificate: "ANIMALQUARANTINECERTIFICATE",
	FeeTypeSpecialWarehousingFee:       "SPECIALWAREHOUSINGFEE",
	FeeTypeSurcharge:                   "SURCHARGE",
	FeeTypeExpressCompensation:         "EXPRESSCOMPENSATION",
	FeeTypeReceiveDeliveryFeeClaims:    "RECEIVEDELIVERYFEECLAIMS",
	FeeTypeExpressAddrChangeFee:        "EXPRESSADDRCHANGEFEE",
	FeeTypeExceededAreaFee:             "EXCEEDEDAREAFEE",
	FeeTypeOverLengthOverWeightFee:     "OVERLENGTHOVERWEIGHTFEE",
	FeeTypeReturnFee:                   "RETURNFEE",
	FeeTypeFreshFee:                    "FRESHFEE",
	FeeTypeFullInsuranceFee:            "FULLINSURANCEFEE",
	FeeTypeOtherFee:                    "OTHERFEE",
}

var FeeTypeReverseMap = map[string]FeeType{
	"INSURANCEFEE":                FeeTypeInsuranceFee,
	"PACKAGINGFEE":                FeeTypePackagingFee,
	"COLLECTIONFEE":               FeeTypeCollectionFee,
	"SIGNANDRETURN":               FeeTypeSignAndReturn,
	"HANDLINGCHARGES":             FeeTypeHandlingCharges,
	"PICKUPSERVICEFEE":            FeeTypePickUpServiceFee,
	"DELIVERYSERVICEFEE":          FeeTypeDeliveryServiceFee,
	"RESOURCECONDITIONINGCHARGES": FeeTypeResourceConditioningCharges,
	"ANIMALQUARANTINECERTIFICATE": FeeTypeAnimalQuarantineCertificate,
	"SPECIALWAREHOUSINGFEE":       FeeTypeSpecialWarehousingFee,
	"SURCHARGE":                   FeeTypeSurcharge,
	"EXPRESSCOMPENSATION":         FeeTypeExpressCompensation,
	"RECEIVEDELIVERYFEECLAIMS":    FeeTypeReceiveDeliveryFeeClaims,
	"EXPRESSADDRCHANGEFEE":        FeeTypeExpressAddrChangeFee,
	"EXCEEDEDAREAFEE":             FeeTypeExceededAreaFee,
	"OVERLENGTHOVERWEIGHTFEE":     FeeTypeOverLengthOverWeightFee,
	"RETURNFEE":                   FeeTypeReturnFee,
	"FRESHFEE":                    FeeTypeFreshFee,
	"FULLINSURANCEFEE":            FeeTypeFullInsuranceFee,
	"OTHERFEE":                    FeeTypeOtherFee,
}

var FeeTypeDescMap = map[FeeType]string{
	FeeTypeInsuranceFee:                "保价费",
	FeeTypePackagingFee:                "包装费",
	FeeTypeCollectionFee:               "代收货款",
	FeeTypeSignAndReturn:               "签单返还",
	FeeTypeHandlingCharges:             "装卸服务",
	FeeTypePickUpServiceFee:            "提货服务",
	FeeTypeDeliveryServiceFee:          "配送服务",
	FeeTypeResourceConditioningCharges: "资源调节费",
	FeeTypeAnimalQuarantineCertificate: "动检证",
	FeeTypeSpecialWarehousingFee:       "特殊入仓",
	FeeTypeSurcharge:                   "附加费",
	FeeTypeExpressCompensation:         "快递赔付费",
	FeeTypeReceiveDeliveryFeeClaims:    "收派服务费-理赔",
	FeeTypeExpressAddrChangeFee:        "快递改址费",
	FeeTypeExceededAreaFee:             "超区服务费",
	FeeTypeOverLengthOverWeightFee:     "超长超重附加费",
	FeeTypeReturnFee:                   "逆向费用",
	FeeTypeFreshFee:                    "保鲜服务费",
	FeeTypeFullInsuranceFee:            "顺丰足额保",
	FeeTypeOtherFee:                    "其他费用",
}

var FeeTypeReverseDescMap = map[string]string{
	"INSURANCEFEE":                "保价费",
	"PACKAGINGFEE":                "包装费",
	"COLLECTIONFEE":               "代收货款",
	"SIGNANDRETURN":               "签单返还",
	"HANDLINGCHARGES":             "装卸服务",
	"PICKUPSERVICEFEE":            "提货服务",
	"DELIVERYSERVICEFEE":          "配送服务",
	"RESOURCECONDITIONINGCHARGES": "资源调节费",
	"ANIMALQUARANTINECERTIFICATE": "动检证",
	"SPECIALWAREHOUSINGFEE":       "特殊入仓",
	"SURCHARGE":                   "附加费",
	"EXPRESSCOMPENSATION":         "快递赔付费",
	"RECEIVEDELIVERYFEECLAIMS":    "收派服务费-理赔",
	"EXPRESSADDRCHANGEFEE":        "快递改址费",
	"EXCEEDEDAREAFEE":             "超区服务费",
	"OVERLENGTHOVERWEIGHTFEE":     "超长超重附加费",
	"RETURNFEE":                   "逆向费用",
	"FRESHFEE":                    "保鲜服务费",
	"FULLINSURANCEFEE":            "顺丰足额保",
	"OTHERFEE":                    "其他费用",
}

type PayStatus int8

const (
	PayFail      = -1 // 支付失败
	PayUnpaid    = 0  // 未支付
	PayPaid      = 1  // 已支付
	PayNoPayment = 2  // 无需支付
	PayRefund    = 3  // 已退款
)

// CompanyCode 快递公司编码
type CompanyCode string

const (
	CompanyCodeYTO  CompanyCode = "yuantong"         // 圆通
	CompanyCodeYD   CompanyCode = "yunda"            // 韵达
	CompanyCodeZTO  CompanyCode = "zhongtong"        // 中通
	CompanyCodeSTO  CompanyCode = "shentong"         // 申通
	CompanyCodeJTEX CompanyCode = "jtexpress"        // 极兔速递
	CompanyCodeSF   CompanyCode = "shunfeng"         // 顺丰速运
	CompanyCodeEMS  CompanyCode = "ems"              // EMS
	CompanyCodeYZPY CompanyCode = "youzhengguonei"   // 邮政包裹
	CompanyCodeJD   CompanyCode = "jd"               // 京东
	CompanyCodeYZBK CompanyCode = "youzhengbk"       // 邮政标准快递
	CompanyCodeDBL  CompanyCode = "debangkuaidi"     // 德邦快递
	CompanyCodeSFKY CompanyCode = "shunfengkuaiyun"  // 顺丰快运
	CompanyCodeDN   CompanyCode = "danniao"          // 丹鸟
	CompanyCodeDB   CompanyCode = "debangwuliu"      // 德邦物流
	CompanyCodeJDKY CompanyCode = "jingdongkuaiyun"  // 京东快运
	CompanyCodeZTKY CompanyCode = "zhongtongkuaiyun" // 中通快运
	CompanyCodeHTKY CompanyCode = "huitongkuaidi"    // 百世快递
	CompanyCodeFWSY CompanyCode = "fengwang"         // 丰网速运
)

// kuaidi100ToStandardMapper 快递100转标准快递公司编码
var kuaidi100ToStandardMapper = map[CompanyCode]express.StandardCourierCode{
	CompanyCodeYTO:  express.StandardCourierCodeYTO,
	CompanyCodeYD:   express.StandardCourierCodeYD,
	CompanyCodeZTO:  express.StandardCourierCodeZTO,
	CompanyCodeSTO:  express.StandardCourierCodeSTO,
	CompanyCodeJTEX: express.StandardCourierCodeJTEX,
	CompanyCodeSF:   express.StandardCourierCodeSF,
	CompanyCodeEMS:  express.StandardCourierCodeEMS,
	CompanyCodeYZPY: express.StandardCourierCodeYZPY,
	CompanyCodeJD:   express.StandardCourierCodeJD,
	CompanyCodeYZBK: express.StandardCourierCodeYZBK,
	CompanyCodeDBL:  express.StandardCourierCodeDBL,
	CompanyCodeSFKY: express.StandardCourierCodeSFKY,
	CompanyCodeDN:   express.StandardCourierCodeDN,
	CompanyCodeDB:   express.StandardCourierCodeDB,
	CompanyCodeJDKY: express.StandardCourierCodeJDKY,
	CompanyCodeZTKY: express.StandardCourierCodeZTKY,
	CompanyCodeHTKY: express.StandardCourierCodeHTKY,
	CompanyCodeFWSY: express.StandardCourierCodeFWSY,
}

// GetStandardCode 获取标准快递公司编码
func (e CompanyCode) GetStandardCode() express.StandardCourierCode {
	return kuaidi100ToStandardMapper[e]
}

// GetStandardCode 获取标准快递公司编码
func GetStandardCode(companyCode string) express.StandardCourierCode {
	return kuaidi100ToStandardMapper[CompanyCode(companyCode)]
}

// StandardToKuaidi100Mapper 标准快递公司编码转快递100
var StandardToKuaidi100Mapper = map[express.StandardCourierCode]CompanyCode{
	express.StandardCourierCodeYTO:  CompanyCodeYTO,
	express.StandardCourierCodeYD:   CompanyCodeYD,
	express.StandardCourierCodeZTO:  CompanyCodeZTO,
	express.StandardCourierCodeSTO:  CompanyCodeSTO,
	express.StandardCourierCodeJTEX: CompanyCodeJTEX,
	express.StandardCourierCodeSF:   CompanyCodeSF,
	express.StandardCourierCodeEMS:  CompanyCodeEMS,
	express.StandardCourierCodeYZPY: CompanyCodeYZPY,
	express.StandardCourierCodeJD:   CompanyCodeJD,
	express.StandardCourierCodeYZBK: CompanyCodeYZBK,
	express.StandardCourierCodeDBL:  CompanyCodeDBL,
	express.StandardCourierCodeSFKY: CompanyCodeSFKY,
	express.StandardCourierCodeDN:   CompanyCodeDN,
	express.StandardCourierCodeDB:   CompanyCodeDB,
	express.StandardCourierCodeJDKY: CompanyCodeJDKY,
	express.StandardCourierCodeZTKY: CompanyCodeZTKY,
	express.StandardCourierCodeHTKY: CompanyCodeHTKY,
	express.StandardCourierCodeFWSY: CompanyCodeFWSY,
}

// GetStandardMapCompanyCode 获取标准快递公司编码对应的快递100编码
func GetStandardMapCompanyCode(standardCode express.StandardCourierCode) CompanyCode {
	return StandardToKuaidi100Mapper[standardCode]
}

// Kuaidi100CompanyCodeNameMapper 快递100公司编码对应的名称
var Kuaidi100CompanyCodeNameMapper = map[CompanyCode]string{
	CompanyCodeYTO:  "圆通",
	CompanyCodeYD:   "韵达",
	CompanyCodeZTO:  "中通",
	CompanyCodeSTO:  "申通",
	CompanyCodeJTEX: "极兔速递",
	CompanyCodeSF:   "顺丰速运",
	CompanyCodeEMS:  "EMS",
	CompanyCodeYZPY: "邮政包裹",
	CompanyCodeJD:   "京东",
	CompanyCodeYZBK: "邮政标准快递",
	CompanyCodeDBL:  "德邦快递",
	CompanyCodeSFKY: "顺丰快运",
	CompanyCodeDN:   "丹鸟",
	CompanyCodeDB:   "德邦物流",
	CompanyCodeJDKY: "京东快运",
	CompanyCodeZTKY: "中通快运",
	CompanyCodeHTKY: "百世快递",
	CompanyCodeFWSY: "丰网速运",
}

// Name 获取快递公司名称
func (e CompanyCode) Name() string {
	return Kuaidi100CompanyCodeNameMapper[e]
}

type ServiceType string

const (
	ServiceTypeTHS      ServiceType = "特惠送" // 特惠送
	ServiceTypeDBDJ360              = "德邦大件360"
	ServiceTypeSFBK                 = "顺丰标快"
	ServiceTypeSFTK                 = "顺丰特快"
	ServiceTypeStandard             = "标准快递"
)

var ServiceTypeMap = map[CompanyCode]ServiceType{
	CompanyCodeJD:   ServiceTypeTHS,
	CompanyCodeDB:   ServiceTypeDBDJ360,
	CompanyCodeDBL:  ServiceTypeStandard,
	CompanyCodeSF:   ServiceTypeSFBK,
	CompanyCodeSFKY: ServiceTypeSFTK,
}

// GetCompanyCodeServiceType 获取快递公司编码对应的服务类型
func GetCompanyCodeServiceType(compCode CompanyCode) ServiceType {
	val, ok := ServiceTypeMap[compCode]
	if !ok {
		return ServiceTypeStandard
	}
	return val
}
