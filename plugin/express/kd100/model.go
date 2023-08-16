package kd100

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
	CallbackUrl string `json:"callbackurl"`    // 回调地址
	Salt        string `json:"salt,omitempty"` // MD5加密串

	// 添加此字段表示开通行政区域解析功能。空：关闭（默认），
	// 3：开通行政区域解析功能以及物流轨迹增加物流状态名称
	// 5: 开通行政解析功能以及物流轨迹增加物流高级状态名称、状态值并且返回出发、目的及当前城市信息
	Resultv2           string `json:"resultv2"`
	AutoCom            string `json:"autoCom"`            // 开通智能单号识别功能：0-关闭（默认），1-开通
	InterCom           string `json:"interCom"`           // 开通国际版功能：0-关闭（默认），1-开通
	DepartureCountry   string `json:"departureCountry"`   // 出发国
	DepartureCom       string `json:"departureCom"`       // 出发国快递公司编码
	DestinationCountry string `json:"destinationCountry"` // 目的国
	DestinationCom     string `json:"destinationCom"`     // 目的国快递公司编码
	Phone              string `json:"phone"`              // 手机号 顺丰速运、顺丰快运和丰网速运必填
}

type SubscribeMapTrackParam struct {
	Key        string                       `json:"key"`        // 授权key
	Company    string                       `json:"company"`    // 快递公司编码
	Number     string                       `json:"number"`     // 快递单号
	From       string                       `json:"from"`       // 出发地城市
	To         string                       `json:"to"`         // 目的地城市
	Parameters *SubscribeMapTrackParameters `json:"parameters"` // 可选参数
}

type SubscribeMapTrackParameters struct {
	CallbackUrl  string `json:"callbackurl"`  // 回调地址
	Salt         string `json:"salt"`         // MD5加密串
	Phone        string `json:"phone"`        // 手机号 顺丰速运、顺丰快运和丰网速运必填
	MapConfigKey string `json:"mapConfigKey"` // 地图轨迹模版ID

	// 添加此字段表示开通行政区域解析功能。空：关闭（默认），
	// 3：开通行政区域解析功能以及物流轨迹增加物流状态名称
	// 5: 开通行政解析功能以及物流轨迹增加物流高级状态名称、状态值并且返回出发、目的及当前城市信息
	ResultV2 string `json:"resultv2"`
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
	Com       string   `json:"com"`       // 快递公司编码
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

type SubscribeMapTrackPush struct {
	Status     string `json:"status"`
	Billstatus string `json:"billstatus"`
	Message    string `json:"message"`
	LastResult struct {
		Message   string `json:"message"`
		State     string `json:"state"`
		Status    string `json:"status"`
		Condition string `json:"condition"`
		Ischeck   string `json:"ischeck"`
		Com       string `json:"com"`
		Nu        string `json:"nu"`
		Data      []struct {
			Context    string `json:"context"`
			Time       string `json:"time"`
			Ftime      string `json:"ftime"`
			Status     string `json:"status"`
			AreaCode   string `json:"areaCode"`
			AreaName   string `json:"areaName"`
			Location   string `json:"location"`
			AreaCenter string `json:"areaCenter"`
			AreaPinYin string `json:"areaPinYin"`
			StatusCode string `json:"statusCode"`
		} `json:"data"`
		RouteInfo struct {
			From struct {
				Number string `json:"number"`
				Name   string `json:"name"`
			} `json:"from"`
			Cur struct {
				Number string `json:"number"`
				Name   string `json:"name"`
			} `json:"cur"`
			To struct {
				Number string `json:"number"`
				Name   string `json:"name"`
			} `json:"to"`
		} `json:"routeInfo"`
		IsLoop      bool   `json:"isLoop"`
		TrailUrl    string `json:"trailUrl"`
		ArrivalTime string `json:"arrivalTime"`
		TotalTime   string `json:"totalTime"`
		RemainTime  string `json:"remainTime"`
	} `json:"lastResult"`
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

type QueryMapTrackParams struct {
	QueryParams
	MapConfigKey string `json:"mapConfigKey"` // 地图轨迹模版ID
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
		AreaCenter string `json:"areaCenter"` // 物流事件发生城市经纬度
		AreaPinYin string `json:"areaPinYin"` // 物流事件发生城市的拼音
		StatusCode string `json:"statusCode"` // 物流事件的状态码
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

type QueryMapTrackResponse struct {
	Message     string                          `json:"message"`     // 消息体, 可以忽略
	State       string                          `json:"state"`       // 快递单当前状态，默认为0在途，1揽收，2疑难，3签收，4退签，5派件，8清关，14拒签等10个基础物流状态，如需要返回高级物流状态，请参考 resultv2 传值
	Status      string                          `json:"status"`      // 通讯状态, 可以忽略
	Condition   string                          `json:"condition"`   // 快递单明细状态标记，暂未实现，可以忽略
	IsCheck     int                             `json:"ischeck"`     // 是否签收标记 0:未签收 1:已签收
	Com         string                          `json:"com"`         // 快递公司编码
	Nu          string                          `json:"nu"`          // 快递单号
	TrailUrl    string                          `json:"trailUrl"`    // 轨迹地图链接
	ArrivalTime string                          `json:"arrivalTime"` // 预计到达时间
	TotalTime   string                          `json:"totalTime"`   // 平均耗时
	RemainTime  string                          `json:"remainTime"`  // 剩余时间
	IsLoop      bool                            `json:"isLoop"`      // 是否循环推送
	RouteInfo   *QueryMapTrackResponseRouteInfo `json:"routeInfo"`   // 路由信息
	Data        []*QueryMapTrackResponseData    `json:"data"`        // 轨迹数据
}

type QueryMapTrackResponseRouteInfo struct {
	From struct { // 出发地
		Number string `json:"number"` // 出发地编码
		Name   string `json:"name"`   // 出发地名称
	} `json:"from"`
	Cur struct { // 当前所在地
		Number string `json:"number"` // 当前所在地编码
		Name   string `json:"name"`   // 当前所在地名称
	} `json:"cur"`
	To struct {
		Number string `json:"number"` // 目的地编码
		Name   string `json:"name"`   // 目的地名称
	} `json:"to"`
}

type QueryMapTrackResponseData struct {
	Context    string `json:"context"`    // 物流事件的描述
	Time       string `json:"time"`       // 物流事件发生的时间, 原始格式
	Ftime      string `json:"ftime"`      // 物流事件发生的时间, 格式化后的
	Status     string `json:"status"`     // 物流事件的状态
	AreaCode   string `json:"areaCode"`   // 物流事件发生的区域编码
	AreaName   string `json:"areaName"`   // 物流事件发生的区域名称
	StatusCode string `json:"statusCode"` // 物流事件的状态码
	AreaCenter string `json:"areaCenter"` // 物流事件发生城市经纬度
	Location   string `json:"location"`   // 物流事件的城市
	AreaPinYin string `json:"areaPinYin"` // 物流事件发生城市的拼音
}

// --------------------------- 响应 ---------------------------
// =========================== 快递查询 ===========================

// =========================== 寄件 ===========================
// --------------------------- 商户寄件 ---------------------------

// BCreateShipmentOrderReq 寄件订单创建请求
type BCreateShipmentOrderReq struct {
	Kuaidicom        string `json:"kuaidicom"`        // 快递公司编码
	RecManName       string `json:"recManName"`       // 收件人姓名
	RecManMobile     string `json:"recManMobile"`     // 收件人手机号
	RecManPrintAddr  string `json:"recManPrintAddr"`  // 收件人打印地址
	SendManName      string `json:"sendManName"`      // 寄件人姓名
	SendManMobile    string `json:"sendManMobile"`    // 寄件人手机号
	SendManPrintAddr string `json:"sendManPrintAddr"` // 寄件人打印地址
	Cargo            string `json:"cargo"`            // 货物名称
	CallBackUrl      string `json:"callBackUrl"`      // 回调地址
	Payment          string `json:"payment"`          // 代收货款金额
	ServiceType      string `json:"serviceType"`      // 服务类型
	Weight           string `json:"weight"`           // 货物重量
	Remark           string `json:"remark"`           // 备注

	// 签名用随机字符串，用于验证签名sign。
	//salt值不为null时，推送数据将包含该加密签名，加密方式：md5(param+salt)。
	//注意： salt值为空串时，推送的数据也会包含sign，此时可忽略sign的校验。
	Salt            string `json:"salt"`
	DayType         string `json:"dayType"`         // 时效类型
	PickupStartTime string `json:"pickupStartTime"` // 上门取件开始时间
	PickupEndTime   string `json:"pickupEndTime"`   // 上门取件结束时间
	PasswordSigning string `json:"passwordSigning"` // 密码签收
	ValinsPay       string `json:"valinsPay"`       // 保价额度，单位：元
	Op              string `json:"op"`              // 是否开启订阅功能 0：不开启(默认) 1：开启 说明开启订阅功能时：pollCallBackUrl必须填入 此功能只针对有快递单号的单
	PollCallBackUrl string `json:"pollCallBackUrl"` // 如果op设置为1时，pollCallBackUrl必须填入，用于跟踪回调，回调内容通过五、快递信息推送接口返回（免费服务）

	// 添加此字段表示开通行政区域解析功能 。
	//0：关闭（默认）
	//1：开通行政区域解析功能以及物流轨迹增加物流状态名称 (详见：快递信息推送接口文档)
	//3：开通行政区域解析功能以及物流轨迹增加物流状态名称，同时返回地图内容(详见：地图轨迹推送接口文档)
	Resultv2         string `json:"resultv2"`
	ReturnType       string `json:"returnType"`       // 面单返回类型，默认为空，不返回面单内容。10：设备打印，20：生成图片短链回调。
	Siid             string `json:"siid"`             // 设备码，returnType为10时必填
	Tempid           string `json:"tempid"`           // 模板编码，通过管理后台的电子面单模板信息获取 ，returnType不为空时必填
	PrintCallBackUrl string `json:"printCallBackUrl"` // 打印状态回调地址，returnType为10时必填
}

// BCreateShipmentOrderResp 寄件订单创建响应
type BCreateShipmentOrderResp struct {
	Result     bool   `json:"result"`     // 接口调用结果
	ReturnCode string `json:"returnCode"` // 接口调用结果编码
	Message    string `json:"message"`    // 接口调用结果描述
	Data       struct {
		TaskId    string `json:"taskId"`    // 任务ID
		OrderId   string `json:"orderId"`   // 订单号
		Kuaidinum string `json:"kuaidinum"` // 快递单号
		EOrder    string `json:"eOrder"`    // 电子面单号
	} `json:"data"` // 接口调用结果数据
}

// CreateShipmentOrderPush 寄件订单回调推送数据
type CreateShipmentOrderPush struct {
	Kuaidicom string `json:"kuaidicom"` // 快递公司编码
	Kuaidinum string `json:"kuaidinum"` // 快递单号
	Status    string `json:"status"`    // 订单状态
	Message   string `json:"message"`   // 订单状态描述
	Data      struct {
		OrderId       string `json:"orderId"`       // 订单号
		Status        int    `json:"status"`        // 订单状态
		CourierName   string `json:"courierName"`   // 快递员姓名
		CourierMobile string `json:"courierMobile"` // 快递员手机号
		Weight        string `json:"weight"`        // 货物重量
		DefPrice      string `json:"defPrice"`      // 货物保价金额
		Freight       string `json:"freight"`       // 运费
		Volume        string `json:"volume"`        // 体积
		ActualWeight  string `json:"actualWeight"`  // 实际重量
		FeeDetails    []struct {
			FeeType   string `json:"feeType"`   // 费用类型
			FeeDesc   string `json:"feeDesc"`   // 费用描述
			Amount    string `json:"amount"`    // 费用金额
			PayStatus string `json:"payStatus"` // 支付状态
		} `json:"feeDetails"` // 费用明细
		PrintTaskId string `json:"printTaskId"` // 打印任务ID
		ImgBase64   string `json:"imgBase64"`   // 面单图片base64编码
	} `json:"data"` // 接口调用结果数据
}

// BCancelShipmentOrderReq 寄件订单取消请求
type BCancelShipmentOrderReq struct {
	TaskId    string `json:"taskId"`    // 任务ID
	OrderId   string `json:"orderId"`   // 订单号
	CancelMsg string `json:"cancelMsg"` // 取消原因
}

// BCancelShipmentOrderResp 寄件订单取消响应
type BCancelShipmentOrderResp struct {
	Result     bool   `json:"result"`     // 接口调用结果
	ReturnCode string `json:"returnCode"` // 接口调用结果编码
	Message    string `json:"message"`    // 接口调用结果描述
}

// BSendExpressOrderPush 寄件订单推送数据
type BSendExpressOrderPush struct {
	// 监控状态:polling:监控中，shutdown:结束，abort:中止，updateall：重新推送。
	// 其中当快递单为已签收时status=shutdown，当message为“3天查询无记录”或“60天无变化时”status= abort ，
	// 对于status=abort的状态，需要增加额外的处理逻辑
	Status     string `json:"status"`
	Billstatus string `json:"billstatus"` // 包括got、sending、check三个状态，由于意义不大，已弃用，请忽略
	Message    string `json:"message"`    // 监控状态相关消息，如:3天查询无记录，60天无变化
	LastResult struct {
		Message string `json:"message"`
		Nu      string `json:"nu"`
		Ischeck string `json:"ischeck"`
		Com     string `json:"com"`
		Status  string `json:"status"`
		Data    []struct {
			Time       string `json:"time"`
			Context    string `json:"context"`
			Ftime      string `json:"ftime"`
			AreaCode   string `json:"areaCode"`
			AreaName   string `json:"areaName"`
			Status     string `json:"status"`
			Location   string `json:"location"`
			AreaCenter string `json:"areaCenter"`
			AreaPinYin string `json:"areaPinYin"`
			StatusCode string `json:"statusCode"`
		} `json:"data"`
		State     string `json:"state"`
		Condition string `json:"condition"`
		RouteInfo struct {
			From struct {
				Number string `json:"number"`
				Name   string `json:"name"`
			} `json:"from"`
			Cur struct {
				Number string `json:"number"`
				Name   string `json:"name"`
			} `json:"cur"`
			To interface{} `json:"to"`
		} `json:"routeInfo"`
		IsLoop bool `json:"isLoop"`
	} `json:"lastResult"`
}

// BQueryShipmentOrderReq 寄件订单查询请求
type BQueryShipmentOrderReq struct {
	TaskId string `json:"taskId"` // 任务ID
}

// BQueryShipmentOrderResp 寄件订单查询响应
type BQueryShipmentOrderResp struct {
	Data struct {
		Cargo         string `json:"cargo"`
		Comment       string `json:"comment"`
		CourierMobile string `json:"courierMobile"`
		CourierName   string `json:"courierName"`
		CreateTime    string `json:"createTime"`
		DayType       string `json:"dayType"`
		Freight       string `json:"freight"`
		FeeDetails    []struct {
			FeeType   string `json:"feeType"`
			FeeDesc   string `json:"feeDesc"`
			Amount    string `json:"amount"`
			PayStatus string `json:"payStatus"`
		} `json:"feeDetails"`
		KuaidiCom       string `json:"kuaidiCom"`
		KuaidiNum       string `json:"kuaidiNum"`
		LastWeight      string `json:"lastWeight"`
		OrderId         string `json:"orderId"`
		PayStatus       int    `json:"payStatus"`
		Payment         string `json:"payment"`
		PickupEndTime   string `json:"pickupEndTime"`
		PickupStartTime string `json:"pickupStartTime"`
		PreWeight       string `json:"preWeight"`
		DefPrice        string `json:"defPrice"`
		RecAddr         string `json:"recAddr"`
		RecCity         string `json:"recCity"`
		RecDistrict     string `json:"recDistrict"`
		RecMobile       string `json:"recMobile"`
		RecName         string `json:"recName"`
		RecProvince     string `json:"recProvince"`
		SendAddr        string `json:"sendAddr"`
		SendCity        string `json:"sendCity"`
		SendDistrict    string `json:"sendDistrict"`
		SendMobile      string `json:"sendMobile"`
		SendName        string `json:"sendName"`
		SendProvince    string `json:"sendProvince"`
		ServiceType     string `json:"serviceType"`
		Status          int    `json:"status"`
		TaskId          string `json:"taskId"`
		Valins          string `json:"valins"`
	} `json:"data"`
	Message    string `json:"message"`
	Result     bool   `json:"result"`
	ReturnCode string `json:"returnCode"`
}

// QueryShipmentPriceReq 查询寄件价格请求
type QueryShipmentPriceReq struct {
	Kuaidicom        string `json:"kuaidicom"`
	SendManPrintAddr string `json:"sendManPrintAddr"`
	RecManPrintAddr  string `json:"recManPrintAddr"`
	Weight           string `json:"weight"`
	ServiceType      string `json:"serviceType"`
}

type QueryShipmentPriceResp struct {
	Data struct {
		FirstPrice    string `json:"firstPrice"`
		DefPrice      string `json:"defPrice"`
		DefFirstPrice string `json:"defFirstPrice"`
		Price         string `json:"price"`
		ServiceType   string `json:"serviceType"`
		OverPrice     string `json:"overPrice"`
		DefOverPrice  string `json:"defOverPrice"`
		KuaidiCom     string `json:"kuaidiCom"`
	} `json:"data"`
	Message    string `json:"message"`
	Result     bool   `json:"result"`
	ReturnCode string `json:"returnCode"`
}

// --------------------------- 商户寄件 ---------------------------

// --------------------------- C端寄件 ---------------------------

// --------------------------- C端寄件 ---------------------------
// =========================== 寄件 ===========================
