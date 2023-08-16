package express

// StandardCourierCode // 统一标准快递公司编码
type StandardCourierCode string

const (
	StandardCourierCodeYTO  = "YTO"  // 圆通
	StandardCourierCodeYD   = "YD"   // 韵达
	StandardCourierCodeZTO  = "ZTO"  // 中通
	StandardCourierCodeSTO  = "STO"  // 申通
	StandardCourierCodeJTEX = "JTEX" // 极兔速递
	StandardCourierCodeSF   = "SF"   // 顺丰速运
	StandardCourierCodeEMS  = "EMS"  // EMS
	StandardCourierCodeYZPY = "YZPY" // 邮政包裹
	StandardCourierCodeJD   = "JD"   // 京东
	StandardCourierCodeYZBK = "YZBK" // 邮政标准快
	StandardCourierCodeDBL  = "DBL"  // 德邦快递
	StandardCourierCodeSFKY = "SFKY" // 顺丰快运
	StandardCourierCodeDN   = "DN"   // 丹鸟
	StandardCourierCodeDB   = "DB"   // 德邦物流
	StandardCourierCodeJDKY = "JDKY" // 京东快运
	StandardCourierCodeZTKY = "ZTKY" // 中通快运
	StandardCourierCodeHTKY = "HTKY" // 百世快递
	StandardCourierCodeFWSY = "FWSY" // 丰网速运
)

// StandardCourierCodeNameMapper 统一标准快递公司编码名称映射
var StandardCourierCodeNameMapper = map[StandardCourierCode]string{
	StandardCourierCodeYTO:  "圆通",
	StandardCourierCodeYD:   "韵达",
	StandardCourierCodeZTO:  "中通",
	StandardCourierCodeSTO:  "申通",
	StandardCourierCodeJTEX: "极兔速递",
	StandardCourierCodeSF:   "顺丰速运",
	StandardCourierCodeEMS:  "EMS",
	StandardCourierCodeYZPY: "邮政包裹",
	StandardCourierCodeJD:   "京东",
	StandardCourierCodeYZBK: "邮政标准快递",
	StandardCourierCodeDBL:  "德邦快递",
	StandardCourierCodeSFKY: "顺丰快运",
	StandardCourierCodeDN:   "丹鸟",
	StandardCourierCodeDB:   "德邦物流",
	StandardCourierCodeJDKY: "京东快运",
	StandardCourierCodeZTKY: "中通快运",
	StandardCourierCodeHTKY: "百世快递",
	StandardCourierCodeFWSY: "丰网速运",
}

func (s StandardCourierCode) ZHName() string {
	return StandardCourierCodeNameMapper[s]
}

// State 快递单当前状态
type State uint8

const (
	StateInTransit       = 0  // 在途
	StatePickUp          = 1  // 揽件
	StateDifficult       = 2  // 疑难件
	StateSign            = 3  // 签收
	StateReject          = 4  // 退签
	StateDispatch        = 5  // 派件
	StateReturned        = 6  // 退回
	StateTransferAnother = 7  // 转投
	StateClearance       = 8  // 清关
	StateRefuse          = 14 // 拒签
)

// PayType 付款方式
type PayType uint8

const (
	PayTypeSender    = iota + 1 // 寄付
	PayTypeRecipient            // 到付
	PayTypeMonthly              // 月结
)
