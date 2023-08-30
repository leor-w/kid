package express

import "time"

type CustomParams map[string]interface{}

// ======================== 实时查询 ========================

// QueryParams 查询参数
type QueryParams struct {
	CompCode     StandardCourierCode // 快递公司编码
	ExpNo        string              // 快递单号
	Phone        string              // 快递对应收件/寄件人手机号
	From         string              // 出发地城市
	To           string              // 目的地城市
	Sort         string              // 排序 desc asc
	CustomParams CustomParams        // 自定义参数
}

// QueryResponse 查询响应
type QueryResponse struct {
	State    State               `json:"state"`     // 快递单当前状态
	IsSigned bool                `json:"is_signed"` // 是否签收
	CompCode StandardCourierCode `json:"comp_code"` // 快递公司编码
	ExpNo    string              `json:"exp_no"`    // 快递单号
	Tracks   []*Track            `json:"tracks"`    // 轨迹数据
}

// Track 物流轨迹
type Track struct {
	Description string    `json:"description"` // 物流事件的描述
	Time        time.Time `json:"time"`        // 物流事件发生的时间, 原始格式
	AreaCode    string    `json:"area_code"`   // 物流事件发生的区域编码
	Area        string    `json:"area"`        // 物流事件发生的区域
	Status      State     `json:"status"`      // 物流事件的状态
	Location    string    `json:"location"`    // 物流事件的城市
	Coordinates string    `json:"coordinates"` // 物流事件的城市经纬度坐标
}

type QueryMapResponse struct {
	State       State                   `json:"state"`        // 快递单当前状态
	IsSigned    bool                    `json:"is_signed"`    // 是否签收
	CompCode    StandardCourierCode     `json:"comp_code"`    // 快递公司编码
	ExpNo       string                  `json:"exp_no"`       // 快递单号
	TrailUrl    string                  `json:"trail_url"`    // 轨迹地图链接
	ArrivalTime time.Time               `json:"arrival_time"` // 预计到达时间
	TotalTime   string                  `json:"total_time"`   // 平均耗时
	RemainTime  string                  `json:"remain_time"`  // 剩余时间
	Data        []*QueryMapResponseData `json:"data"`         // 轨迹数据
}

type QueryMapResponseData struct {
	Describe   string    `json:"describe"`   // 物流事件的描述
	Time       time.Time `json:"time"`       // 物流事件发生的时间
	AreaCode   string    `json:"area_code"`  // 物流事件发生的区域编码
	Area       string    `json:"area"`       // 物流事件发生的区域
	State      State     `json:"state"`      // 物流事件的状态码
	Coordinate string    `json:"coordinate"` // 物流事件的城市经纬度坐标
	Location   string    `json:"location"`   // 物流事件的城市
}

// ======================== 实时查询 ========================

// ======================== 轨迹订阅 ========================

// SubscribeParams 订阅参数
type SubscribeParams struct {
	Schema       string              // 请求协议类型 json
	CompCode     StandardCourierCode // 快递公司编码
	ExpNo        string              // 快递单号
	Phone        string              // 快递对应收件/寄件人手机号
	From         string              // 出发地城市
	CallbackUrl  string              // 回调地址
	To           string              // 目的地城市
	Sort         string              // 排序 desc asc
	CustomParams CustomParams        // 自定义参数
}

// SubscribeCallbackPush 订阅回调推送
type SubscribeCallbackPush struct {
	// 监控状态
	Status            string `json:"status"`             // 物流事件的状态
	Message           string `json:"message"`            // 物流事件的描述
	ShouldResubscribe bool   `json:"should_resubscribe"` // 是否需要重新订阅

	// 快递信息
	State    State               `json:"state"`     // 快递单当前状态
	IsSigned bool                `json:"is_signed"` // 是否检查
	CompCode StandardCourierCode `json:"comp_code"` // 快递公司编码
	ExpNo    string              `json:"exp_no"`    // 快递单号
	Tracks   []*Track            `json:"tracks"`    // 轨迹数据
}

type SubscribeCallbackMapPush struct {
	// 监控状态
	Status            string `json:"status"`             // 物流事件的状态
	Message           string `json:"message"`            // 物流事件的描述
	ShouldResubscribe bool   `json:"should_resubscribe"` // 是否需要重新订阅

	// 快递信息
	State       State               `json:"state"`        // 快递单当前状态
	IsSigned    bool                `json:"is_signed"`    // 是否检查
	CompCode    StandardCourierCode `json:"comp_code"`    // 快递公司编码
	ExpNo       string              `json:"exp_no"`       // 快递单号
	TrailUrl    string              `json:"trail_url"`    // 轨迹地图链接
	ArrivalTime time.Time           `json:"arrival_time"` // 预计到达时间
	TotalTime   string              `json:"total_time"`   // 平均耗时
	RemainTime  string              `json:"remain_time"`  // 剩余时间
	Tracks      []*Track            `json:"tracks"`       // 轨迹数据
}

// ======================== 轨迹订阅 ========================

// ======================== 寄件 ========================

// SendExpressParams 寄件参数
type SendExpressParams struct {
	CompCode         StandardCourierCode // 快递公司编码
	Recipient        string              // 收件人姓名
	RecipientPhone   string              // 收件人手机号
	RecipientAddr    string              // 收件人地址
	SenderName       string              // 寄件人姓名
	SenderPhone      string              // 寄件人手机号
	SenderAddr       string              // 寄件人地址
	ItemName         string              // 物品名称
	Callback         bool                // 是否回调
	CallbackUrl      string              // 回调地址
	PayType          PayType             // 付款方式 1-寄付 2-到付 3-第三方付
	ServiceType      string              // 服务类型 1-标准快递 2-同城快递 3-当日快递 4-次晨快递 5-即时快递 6-国际快递
	Weight           float64             // 物品重量
	Remark           string              // 备注
	DayType          string              // 预约时间
	StartTime        string              // 预约开始时间
	EndTime          string              // 预约结束时间
	InsuredAmount    float64             // 保价金额
	RealName         string              // 实名认证姓名
	SenderIdCardType uint8               // 寄件人证件类型
	SenderIdCard     string              // 寄件人证件号
	PwdSigning       string              // 密码签收
	TradeNo          string              // 平台交易号
	CustomParams     CustomParams        // 自定义参数
}

// SendExpressResponse 寄件响应
type SendExpressResponse struct {
	TaskID  string // 任务ID
	OrderNo string // 订单号
	ExpNo   string // 快递单号
}

type BSendExpressCallbackPush struct {
	CompCode         StandardCourierCode // 快递公司编码
	ExpNo            string              // 快递单号
	OrderId          string              // 订单号
	Status           int                 // 状态
	PickupAgentName  string              // 取件员姓名
	PickupAgentPhone string              // 取件员电话
	Weight           float64             // 物品重量
	Freight          float64             // 运费
	Volume           float64             // 体积
	ActualWeight     float64             // 实际重量
	Fees             []*FeeDetail        // 费用明细
}

type FeeDetail struct {
	FeeType   string // 费用类型
	FeeDesc   string // 费用描述
	Fee       string // 费用
	PayStatus string // 支付状态
}

// CSendExpressCallbackPush C端寄件回调推送
type CSendExpressCallbackPush struct {
	CompCode         StandardCourierCode // 快递公司编码
	ExpNo            string              // 快递单号
	OrderId          string              // 订单号
	Status           int                 // 状态
	PickupAgentName  string              // 取件员姓名
	PickupAgentPhone string              // 取件员电话
	Weight           float64             // 物品重量
	Freight          float64             // 运费
}

// CancelExpressParams 取消寄件参数
type CancelExpressParams struct {
	TaskId  string // 任务ID
	OrderId string // 订单号
	Reason  string // 取消原因
}

// CancelExpressResponse 取消寄件响应
type CancelExpressResponse struct {
	Success bool   // 是否提交成功
	Code    string // 提交结果编码
	Msg     string // 提交结果描述
}

// QueryPostageParams 查询运费参数
type QueryPostageParams struct {
	CompCode     StandardCourierCode // 快递公司编码
	SenderAddr   string              // 寄件人地址
	ReceiverAddr string              // 收件人地址
	Weight       float64             // 物品重量
	ExpType      uint8               // 快递类型
}

type QueryBSendExpressResponse struct {
	TaskId               string              // 任务ID
	CreatedAt            time.Time           // 下单时间
	OrderId              string              // 订单号
	CompCode             StandardCourierCode // 快递公司编码
	ExpNo                string              // 快递单号
	SenderName           string              // 寄件人姓名
	SenderPhone          string              // 寄件人手机号
	SenderProvince       string              // 寄件人省份
	SenderCity           string              // 寄件人城市
	SenderDistrict       string              // 寄件人区县
	SenderAddress        string              // 寄件人详细地址
	ReceiverName         string              // 收件人姓名
	ReceiverPhone        string              // 收件人手机号
	ReceiverProvince     string              // 收件人省份
	ReceiverCity         string              // 收件人城市
	ReceiverDistrict     string              // 收件人区县
	ReceiverAddress      string              // 收件人详细地址
	Valins               string              // 保价金额
	ItemName             string              // 物品名称
	ExpressType          string              // 快递类型
	Weight               float64             // 物品重量
	ActualWeight         float64             // 实际重量
	Remark               string              // 备注
	ReservationData      string              // 预约日期
	ReservationStartTime string              // 预约起始时间
	ReservationEndTime   string              // 预约截止时间
	PayType              PayType             // 付款方式
	Freight              float64             // 运费
	PickupAgentName      string              // 取件员姓名
	PickupAgentPhone     string              // 取件员电话
	Status               int                 // 状态
	PayStatus            uint8               // 支付状态
	FeeDetails           []*FeeDetail        // 费用明细
}

// QueryPostageResponse 查询运费响应
type QueryPostageResponse struct {
	FirstWeight     float64 // 首重价格
	ContinuedWeight float64 // 续重价格
	TotalPrice      float64 // 总价格
	ExpType         string  // 快递类型
}

// ======================== 寄件 ========================
