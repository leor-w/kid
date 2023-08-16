package express

// =========================== 订阅 ===========================
// --------------------------- 请求 ---------------------------

type SubscribeReqConfig struct {
	Schema             string // 返回的数据格式, json(默认),xml
	Company            string // 快递公司编码
	Number             string // 快递单号
	From               string // 出发地城市
	To                 string // 目的地城市
	CallbackUrl        string // 回调地址
	Salt               string // MD5加密串
	Resultv2           string // 开通行政区域解析功能：0-关闭（默认），1-开通
	AutoCom            string // 开通智能单号识别功能：0-关闭（默认），1-开通
	InterCom           string // 开通国际版功能：0-关闭（默认），1-开通
	DepartureCountry   string // 出发国
	DepartureCom       string // 出发国快递公司编码
	DestinationCountry string // 目的国
	DestinationCom     string // 目的国快递公司编码
	Phone              string // 手机号 顺丰速运、顺丰快运和丰网速运必填
}

// SubscribeReq 订阅请求
type SubscribeReq struct {
	Schema string          `json:"schema"` // 返回的数据格式, json(默认),xml
	Param  *SubscribeParam `json:"param"`  // 请求参数
}

// SubscribeParam 请求参数
type SubscribeParam struct {
	Company    string               `json:"company"`        // 快递公司编码
	Number     string               `json:"number"`         // 快递单号
	From       string               `json:"from,omitempty"` // 出发地城市
	To         string               `json:"to,omitempty"`   // 目的地城市
	Key        string               `json:"key"`            // 授权key
	Parameters *SubscribeParameters `json:"parameters"`     // 可选参数
}

// SubscribeParameters 可选参数
type SubscribeParameters struct {
	CallbackUrl        string `json:"callbackurl"`        // 回调地址
	Salt               string `json:"salt,omitempty"`     // MD5加密串
	Resultv2           string `json:"resultv2"`           // 开通行政区域解析功能：0-关闭（默认），1-开通
	AutoCom            string `json:"autoCom"`            // 开通智能单号识别功能：0-关闭（默认），1-开通
	InterCom           string `json:"interCom"`           // 开通国际版功能：0-关闭（默认），1-开通
	DepartureCountry   string `json:"departureCountry"`   // 出发国
	DepartureCom       string `json:"departureCom"`       // 出发国快递公司编码
	DestinationCountry string `json:"destinationCountry"` // 目的国
	DestinationCom     string `json:"destinationCom"`     // 目的国快递公司编码
	Phone              string `json:"phone"`              // 手机号 顺丰速运、顺丰快运和丰网速运必填
}

type SubscribeResp struct {
	Message    string `json:"message"`
	Result     bool   `json:"result"`
	ReturnCode string `json:"returnCode"`
}

// --------------------------- 请求 ---------------------------

// --------------------------- 回调 ---------------------------

type SubscribePush struct {
	Status     string  `json:"status"`     // 监控状态：polling:监控中，shutdown:结束，abort:中止，updateall：重新推送
	Billstatus string  `json:"billstatus"` // 包括got、sending、check三个状态，由于意义不大，已弃用，请忽略
	Message    string  `json:"message"`    // 监控状态相关消息，如:3天查询无记录，60天无变化
	AutoCheck  string  `json:"autoCheck"`  // 是否开启智能判断功能，当快递单号识别出多个快递公司时，是否开启智能判断功能，默认值0，0:关闭智能判断，1:开启智能判断
	ComOld     string  `json:"comOld"`     // 识别出的快递公司编码
	ComNew     string  `json:"comNew"`     // 智能判断出的快递公司编码
	LastResult *Result `json:"lastResult"` // 最新查询结果，包括：comNew、nuNew、ischeck、condition、status、state、data、message、autoCheck、comOld、nuOld
}

type Result struct {
	Message   string   `json:"message"`   // 监控状态相关消息，如:3天查询无记录，60天无变化
	Nu        string   `json:"nu"`        // 快递单号
	Ischeck   string   `json:"ischeck"`   // 是否签收标记 0:未签收 1:已签收
	Com       ExpCom   `json:"com"`       // 快递公司编码
	Status    string   `json:"status"`    // 快递单当前的状态 ：0：在途，即货物处于运输过程中；1：揽件，货物已由快递公司揽收并且产生了第一条跟踪信息；
	Data      []*Data  `json:"data"`      // 快递单明细状态标记，多行
	State     string   `json:"state"`     // 快递单当前状态，默认为0在途，1揽收，2疑难，3签收，4退签，5派件，8清关，14拒签等10个基础物流状态，如需要返回高级物流状态，请参考 resultv2 传值
	Condition string   `json:"condition"` // 快递单明细状态标记，多行
	RouteInfo struct { // 路由信息
		From struct { // 出发地
			Number string `json:"number"`
			Name   string `json:"name"`
		} `json:"from"`
		Cur struct { // 当前所在地
			Number string `json:"number"`
			Name   string `json:"name"`
		} `json:"cur"`
		To interface{} `json:"to"` // 目的地
	} `json:"routeInfo"`
	IsLoop bool `json:"isLoop"` // 是否循环推送
}

type Data struct {
	Time       string `json:"time"`       // 物流事件发生的时间
	Context    string `json:"context"`    // 物流事件的描述
	Ftime      string `json:"ftime"`      // 物流事件发生的时间
	AreaCode   string `json:"areaCode"`   // 物流事件发生的区域编码
	AreaName   string `json:"areaName"`   // 物流事件发生的区域名称
	Status     string `json:"status"`     // 物流事件的状态
	Location   string `json:"location"`   // 物流事件的城市
	AreaCenter string `json:"areaCenter"` // 物流事件的城市
	AreaPinYin string `json:"areaPinYin"` // 物流事件的城市
	StatusCode string `json:"statusCode"` // 物流事件的状态编码
}

// PushResp 推送响应
type PushResp struct {
	Result     bool   `json:"result"`     // 推送处理结果
	ReturnCode string `json:"returnCode"` // 推送处理结果编码
	Message    string `json:"message"`    // 推送处理结果描述
}

// --------------------------- 回调 ---------------------------
// =========================== 订阅 ===========================

// =========================== 快递查询 ===========================
// --------------------------- 请求 ---------------------------

// QueryReqConfig 查询请求配置
type QueryReqConfig struct {
	Com      string // 快递公司编码
	Num      string // 快递单号
	Phone    string // 手机号
	From     string // 出发地城市
	To       string // 目的地城市
	Resultv2 string // 开通行政区域解析功能：0-关闭（默认），1-开通
	Show     string // 返回数据格式：0-json（默认），1-xml，2-html，3-text
	Order    string // 排序：desc-降序（默认），asc-升序
}

// QueryRequest 查询请求
type QueryRequest struct {
	Customer string       `json:"customer"`
	Sign     string       `json:"sign"`
	Param    *QueryParams `json:"param"`
}

// QueryParams 查询请求参数
type QueryParams struct {
	Com      string `json:"com"`      // 快递公司编码
	Num      string `json:"num"`      // 快递单号
	Phone    string `json:"phone"`    // 手机号
	From     string `json:"from"`     // 出发地城市
	To       string `json:"to"`       // 目的地城市
	Resultv2 string `json:"resultv2"` // 开通行政区域解析功能：0-关闭（默认），1-开通
	Show     string `json:"show"`     // 返回数据格式：0-json（默认），1-xml，2-html，3-text
	Order    string `json:"order"`    // 排序：desc-降序（默认），asc-升序
}

// --------------------------- 请求 ---------------------------
// --------------------------- 响应 ---------------------------

// QueryExpressResp 查询响应
type QueryExpressResp struct {
	Message   string     `json:"message"`   // 消息体, 可以忽略
	State     string     `json:"state"`     // 快递单当前状态，默认为0在途，1揽收，2疑难，3签收，4退签，5派件，8清关，14拒签等10个基础物流状态，如需要返回高级物流状态，请参考 resultv2 传值
	Status    string     `json:"status"`    // 通讯状态, 可以忽略
	Condition string     `json:"condition"` // 快递单明细状态标记，暂未实现，可以忽略
	Nu        string     `json:"nu"`        // 快递单号
	Ischeck   string     `json:"ischeck"`   // 是否签收标记 0:未签收 1:已签收
	Com       string     `json:"com"`       // 快递公司编码
	Data      []struct { // 轨迹数据
		Context    string `json:"context"`    // 物流事件的描述
		Time       string `json:"time"`       // 物流事件发生的时间, 原始格式
		Ftime      string `json:"ftime"`      // 物流事件发生的时间, 格式化后的
		AreaCode   string `json:"areaCode"`   // 物流事件发生的区域编码
		AreaName   string `json:"areaName"`   // 物流事件发生的区域名称
		Status     string `json:"status"`     // 物流事件的状态
		Location   string `json:"location"`   // 物流事件的城市
		AreaCenter string `json:"areaCenter"` // 物流事件的城市
		AreaPinYin string `json:"areaPinYin"` // 物流事件的城市
		StatusCode string `json:"statusCode"` // 物流事件的状态编码
	} `json:"data"`
	RouteInfo struct {
		From struct { // 出发地
			Number string `json:"number"` // 出发地编码
			Name   string `json:"name"`   // 出发地名称
		} `json:"from"`
		Cur struct { // 当前所在地
			Number string `json:"number"` // 当前所在地编码
			Name   string `json:"name"`   // 当前所在地名称
		} `json:"cur"`
		To interface{} `json:"to"` // 目的地
	} `json:"routeInfo"`
	IsLoop bool `json:"isLoop"` // 是否循环推送

	// 错误处理
	Result     bool   `json:"result"`     // 查询结果
	ReturnCode string `json:"returnCode"` // 查询结果编码
}

// --------------------------- 响应 ---------------------------
// =========================== 快递查询 ===========================
