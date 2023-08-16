package douyin

// -------------------------------------------- 抖音登录相关 --------------------------------------------

// AccessTokenResponse 抖音获取access_token的返回内容
type AccessTokenResponse struct {
	Data struct {
		AccessToken      string `json:"access_token"`
		Description      string `json:"description"`
		ErrorCode        int    `json:"error_code"`
		ExpiresIn        string `json:"expires_in"`
		OpenId           string `json:"open_id"`
		RefreshExpiresIn string `json:"refresh_expires_in"`
		RefreshToken     string `json:"refresh_token"`
		Scope            string `json:"scope"`
	} `json:"data"`
	Message string `json:"message"`
	Extra   struct {
		Logid string `json:"logid"`
		Now   int64  `json:"now"`
	} `json:"extra"`
}

// UserinfoResponse 抖音获取用户信息的返回内容
type UserinfoResponse struct {
	Data struct {
		Avatar        string `json:"avatar"`
		AvatarLarger  string `json:"avatar_larger"`
		ClientKey     string `json:"client_key"`
		EAccountRole  string `json:"e_account_role"`
		ErrorCode     int    `json:"error_code"`
		LogId         string `json:"log_id"`
		Nickname      string `json:"nickname"`
		OpenId        string `json:"open_id"`
		UnionId       string `json:"union_id"`
		EncryptMobile string `json:"encrypt_mobile"`
		Description   string `json:"description"`
	} `json:"data"`
	Extra struct {
		Logid string `json:"logid"`
		Now   int64  `json:"now"`
	} `json:"extra"`
	Message string `json:"message"`
}

// -------------------------------------------- 抖音登录相关 --------------------------------------------

// -------------------------------------------- 抖音支付相关 --------------------------------------------

// ======================================== 订单创建相关 ========================================

// PayCreateOrderConfig 抖音创建预下单接口配置
type PayCreateOrderConfig struct {
	OutTradeNo      *string
	TotalAmount     *int64
	Subject         *string
	Body            *string
	ValidTime       *int
	CpExtra         *string
	DisableMsg      *int
	MsgPage         *string
	ExpandOrderInfo *ExpandOrderInfo
	LimitPayWay     *string
	NotifyUrl       *string
}

// PayCreateOrderRequest 抖音创建预下单接口请求内容
type PayCreateOrderRequest struct {
	AppId           *string          `json:"app_id"`                      // 小程序APPID
	OutOrderNo      *string          `json:"out_order_no"`                // 开发者侧的订单号。 只能是数字、大小写字母_-*且在同一个app_id下唯一
	TotalAmount     *int64           `json:"total_amount"`                // 支付价格。 单位为[分]
	Subject         *string          `json:"subject"`                     // 商品描述。 长度限制不超过 128 字节且不超过 42 字符
	Body            *string          `json:"body"`                        // 商品详情 长度限制不超过 128 字节且不超过 42 字符
	ValidTime       *int             `json:"valid_time"`                  // 订单过期时间(秒)。最小5分钟，最大2天，小于5分钟会被置为5分钟，大于2天会被置为2天
	Sign            *string          `json:"sign"`                        // 签名
	CpExtra         *string          `json:"cp_extra,omitempty"`          // 开发者自定义字段，回调原样回传。 超过最大长度会被截断
	NotifyUrl       *string          `json:"notify_url,omitempty"`        // 商户自定义回调地址，必须以 https 开头，支持 443 端口。 指定时，支付成功后抖音会请求该地址通知开发者
	ThirdpartyId    *string          `json:"thirdparty_id,omitempty"`     // 第三方平台服务商 id，非服务商模式留空
	StoreUid        *string          `json:"store_uid,omitempty"`         // 可用此字段指定本单使用的收款商户号（目前为灰度功能，需要联系平台运营添加白名单，白名单添加1小时后生效；未在白名单的小程序，该字段不生效）
	DisableMsg      *int             `json:"disable_msg,omitempty"`       // 是否屏蔽支付完成后推送用户抖音消息，1-屏蔽 0-非屏蔽，默认为0。 特别注意： 若接入POI, 请传1。因为POI订单体系会发消息，所以不用再接收一次担保支付推送消息，
	MsgPage         *string          `json:"msg_page,omitempty"`          // 支付完成后推送给用户的抖音消息跳转页面，开发者需要传入在app.json中定义的链接，如果不传则跳转首页。
	ExpandOrderInfo *ExpandOrderInfo `json:"expand_order_info,omitempty"` // 订单拓展信息
	LimitPayWay     *string          `json:"limit_pay_way,omitempty"`     // 屏蔽指定支付方式，屏蔽多个支付方式，请使用逗号","分割，枚举值：
	// 屏蔽微信支付：LIMIT_WX
	// 屏蔽支付宝支付：LIMIT_ALI
	// 屏蔽抖音支付：LIMIT_DYZF
	// 特殊说明：若之前开通了白名单，平台会保留之前屏蔽逻辑；若传入该参数，会优先以传入的为准，白名单则无效
}

type ExpandOrderInfo struct {
	OriginalDeliveryFee *int64 `json:"original_delivery_fee,omitempty"` // 配送费原价，单位为[分]，仅外卖小程序需要传对应信息
	ActualDeliveryFee   *int64 `json:"actual_delivery_fee,omitempty"`   // 实付配送费，单位为[分]，仅外卖小程序需要传对应信息
}

// ToMap 转换为map
func (req *PayCreateOrderRequest) ToMap() map[string]interface{} {
	m := make(map[string]interface{})
	if req.AppId != nil {
		m["app_id"] = *req.AppId
	}
	if req.ThirdpartyId != nil {
		m["thirdparty_id"] = *req.ThirdpartyId
	}
	if req.OutOrderNo != nil {
		m["out_order_no"] = *req.OutOrderNo
	}
	if req.TotalAmount != nil {
		m["total_amount"] = *req.TotalAmount
	}
	if req.Subject != nil {
		m["subject"] = *req.Subject
	}
	if req.Body != nil {
		m["body"] = *req.Body
	}
	if req.ValidTime != nil {
		m["valid_time"] = *req.ValidTime
	}
	if req.NotifyUrl != nil {
		m["notify_url"] = *req.NotifyUrl
	}
	if req.DisableMsg != nil {
		m["disable_msg"] = *req.DisableMsg
	}
	if req.MsgPage != nil {
		m["msg_page"] = *req.MsgPage
	}
	if req.Sign != nil {
		m["sign"] = *req.Sign
	}
	if req.CpExtra != nil {
		m["cp_extra"] = *req.CpExtra
	}
	if req.StoreUid != nil {
		m["store_uid"] = *req.StoreUid
	}
	if req.LimitPayWay != nil {
		m["limit_pay_way"] = *req.LimitPayWay
	}
	return m
}

// PayCreateOrderResponse 抖音创建预下单接口返回内容
type PayCreateOrderResponse struct {
	ErrNo   int    `json:"err_no"`
	ErrTips string `json:"err_tips"`
	Data    struct {
		OrderId    string `json:"order_id"`
		OrderToken string `json:"order_token"`
	} `json:"data"`
}

// ======================================== 订单创建相关 ========================================

// ======================================== 订单回调相关 ========================================

// CallbackReqeust 抖音支付回调请求内容
type CallbackReqeust struct {
	Timestamp    string `json:"timestamp"`     // Unix 时间戳，字符串类型
	Nonce        string `json:"nonce"`         // 随机数
	Msg          string `json:"msg"`           // 订单信息的 json 字符串
	Type         string `json:"type"`          // 回调类型标记，支付成功回调为"payment"
	MsgSignature string `json:"msg_signature"` // 签名
}

// CallbackRequestMessage 抖音支付回调请求订单内容
type CallbackRequestMessage struct {
	Appid          string `json:"appid"`            // 当前交易发起的小程序id
	CpOrderno      string `json:"cp_orderno"`       // 开发者侧的订单号
	CpExtra        string `json:"cp_extra"`         // 预下单时开发者传入字段
	Way            string `json:"way"`              // way 字段中标识了支付渠道： 1-微信支付，2-支付宝支付，10-抖音支付
	ChannelNo      string `json:"channel_no"`       // 支付渠道侧单号(抖音平台请求下游渠道微信或支付宝时传入的单号)
	PaymentOrderNo string `json:"payment_order_no"` // 支付渠道侧PC单号，支付页面可见(微信支付宝侧的订单号)
	TotalAmount    int64  `json:"total_amount"`     // 支付金额，单位为分
	Status         string `json:"status"`           // 固定SUCCESS
	ItemId         string `json:"item_id"`          // 订单来源视频对应视频 id
	SellerUid      string `json:"seller_uid"`       // 该笔交易卖家商户号
	PaidAt         int64  `json:"paid_at"`          // 支付时间，Unix 时间戳，10 位，整型数
	OrderId        string `json:"order_id"`         // 抖音侧订单号
	Extra          string `json:"extra"`            // 该笔交易卖家商户号
}

// CallbackResponse 抖音支付回调返回内容
type CallbackResponse struct {
	ErrNo   int    `json:"err_no"`
	ErrTips string `json:"err_tips"`
}

// ======================================== 订单回调相关 ========================================

// ======================================== 订单查询相关 ========================================

// QueryOrderRequest 抖音支付订单查询请求内容
type QueryOrderRequest struct {
	AppId        *string `json:"app_id"`        // 小程序id
	OutOrderNo   *string `json:"out_order_no"`  // 开发者侧的订单号
	Sign         *string `json:"sign"`          // 签名
	ThirdPartyId *string `json:"thirdparty_id"` // 第三方平台服务商 id，非服务商模式留空
}

func (req *QueryOrderRequest) ToMap() map[string]interface{} {
	m := make(map[string]interface{})
	if req.AppId != nil {
		m["app_id"] = *req.AppId
	}
	if req.OutOrderNo != nil {
		m["out_order_no"] = *req.OutOrderNo
	}
	if req.ThirdPartyId != nil {
		m["thirdparty_id"] = *req.ThirdPartyId
	}
	return m
}

// QueryOrderResponse 抖音支付订单查询返回内容
type QueryOrderResponse struct {
	ErrNo       int    `json:"err_no"`
	ErrTips     string `json:"err_tips"`
	OutOrderNo  string `json:"out_order_no"`
	OrderId     string `json:"order_id"`
	PaymentInfo struct {
		TotalFee    int    `json:"total_fee"`
		OrderStatus string `json:"order_status"`
		PayTime     string `json:"pay_time"`
		Way         int    `json:"way"`
		ChannelNo   string `json:"channel_no"`
		SellerUid   string `json:"seller_uid"`
		ItemId      string `json:"item_id"`
		CpsInfo     string `json:"cps_info"`
	} `json:"payment_info"`
}

const (
	OrderSuccess         = 0    // 订单查询成功
	OrderErrInternal     = 1000 // 内部错误
	OrderErrLimit        = 1001 // 系统限流
	OrderErrNoOrder      = 2000 // 支付记录不存在
	OrderErrSign         = 2008 // 签名错误
	OrderErrParams       = 2010 // 参数错误
	OrderErrAppId        = 2042 // app_id 错误
	OrderErrThirdPartyId = 2047 // thirdparty_id 错误
	OrderErrPermission   = 2048 // 未查询到服务商与小程序的授权关系
)

const (
	OrderStatusSuccess    = "SUCCESS"    // 支付成功
	OrderStatusTimeout    = "TIMEOUT"    // 支付超时
	OrderStatusProcessing = "PROCESSING" // 处理中
	OrderStatusFail       = "FAIL"       // 支付失败
)

// ======================================== 订单查询相关 ========================================

// -------------------------------------------- 抖音支付相关 --------------------------------------------
