package express

// ExpCom 快递公司编码
type ExpCom string

const (
	YTO  = "yuantong"         // 圆通
	YD   = "yunda"            // 韵达
	ZTO  = "zhongtong"        // 中通
	STO  = "shentong"         // 申通
	JTEX = "jtexpress"        // 极兔速递
	SF   = "shunfeng"         // 顺丰速运
	EMS  = "ems"              // EMS
	YZPY = "youzhengguonei"   // 邮政包裹
	JD   = "jd"               // 京东
	YZBK = "youzhengbk"       // 邮政标准快递
	DBL  = "debangkuaidi"     // 德邦快递
	SFKY = "shunfengkuaiyun"  // 顺丰快运
	DN   = "danniao"          // 丹鸟
	DB   = "debangwuliu"      // 德邦物流
	JDKY = "jingdongkuaiyun"  // 京东快运
	ZTKY = "zhongtongkuaiyun" // 中通快运
	HTKY = "huitongkuaidi"    // 百世快递
	FWSY = "fengwang"         // 丰网速运
)

// ExpComAbbrMap 快递公司简称
var ExpComAbbrMap = map[string]string{
	"yuantong":         "YTO",
	"yunda":            "YD",
	"zhongtong":        "ZTO",
	"shentong":         "STO",
	"jtexpress":        "JTEX",
	"shunfeng":         "SF",
	"ems":              "EMS",
	"youzhengguonei":   "YZPY",
	"jd":               "JD",
	"youzhengbk":       "YZBK",
	"debangkuaidi":     "DBL",
	"shunfengkuaiyun":  "SFKY",
	"danniao":          "DN",
	"debangwuliu":      "DB",
	"jingdongkuaiyun":  "JDKY",
	"zhongtongkuaiyun": "ZTKY",
	"huitongkuaidi":    "HTKY",
	"fengwang":         "FWSY",
}

// Abbreviation 获取快递公司简称
func (e ExpCom) Abbreviation() string {
	return ExpComAbbrMap[string(e)]
}

// ExpComFullNameMap 快递公司全称
var ExpComFullNameMap = map[string]string{
	"YTO":  "yuantong",
	"YD":   "yunda",
	"ZTO":  "zhongtong",
	"STO":  "shentong",
	"JTEX": "jtexpress",
	"SF":   "shunfeng",
	"EMS":  "ems",
	"YZPY": "youzhengguonei",
	"JD":   "jd",
	"YZBK": "youzhengbk",
	"DBL":  "debangkuaidi",
	"SFKY": "shunfengkuaiyun",
	"DN":   "danniao",
	"DB":   "debangwuliu",
	"JDKY": "jingdongkuaiyun",
	"ZTKY": "zhongtongkuaiyun",
	"HTKY": "huitongkuaidi",
	"FWSY": "fengwang",
}

// GetExpressCode 获取快递公司编码
func GetExpressCode(abbrName string) string {
	return ExpComFullNameMap[abbrName]
}

// FullName 获取快递公司全称
func (e ExpCom) FullName() string {
	return ExpComFullNameMap[string(e)]
}

// ExpComNameMap 快递公司名称
var ExpComNameMap = map[string]string{
	"yuantong":         "圆通",
	"yunda":            "韵达",
	"zhongtong":        "中通",
	"shentong":         "申通",
	"jtexpress":        "极兔速递",
	"shunfeng":         "顺丰速运",
	"ems":              "EMS",
	"youzhengguonei":   "邮政包裹",
	"jd":               "京东",
	"youzhengbk":       "邮政标准快递",
	"debangkuaidi":     "德邦快递",
	"shunfengkuaiyun":  "顺丰快运",
	"danniao":          "丹鸟",
	"debangwuliu":      "德邦物流",
	"jingdongkuaiyun":  "京东快运",
	"zhongtongkuaiyun": "中通快运",
	"huitongkuaidi":    "百世快递",
	"fengwang":         "丰网速运",
}

// Name 获取快递公司名称
func (e ExpCom) Name() string {
	return ExpComNameMap[string(e)]
}

const (
	ResultTypeClose  = "0" // 关闭,即时查询
	ResultTypeSimple = "1" // 开通行政区域解析功能以及物流轨迹增加物流状态名称
	ResultTypeFull   = "4" // 开通行政解析功能以及物流轨迹增加物流高级状态名称、状态值并且返回出发、目的及当前城市信息
)
