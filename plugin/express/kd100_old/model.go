package kdn

import "net/url"

type CallbackData struct {
	CallbackRequestData
	DataSign    string `json:"data_sign"`
	RequestType string `json:"request_type"`
}

type CallbackRequestData struct {
	PushTime    string `json:"push_time"`
	EBusinessID string `json:"e_business_id"`
	CallbackRequestDataDetail
	Count string `json:"count"`
}

type CallbackRequestDataDetail struct {
	StateEx      string `json:"state_ex"`
	LogisticCode string `json:"logistic_code"`
	ShipperCode  string `json:"shipper_code"`
	Traces       []*Traces
	State        string `json:"state"`
	NextCity     string `json:"next_city"`
	OrderCode    string `json:"order_code"`
	EBusinessID  string `json:"e_business_id"`
	Success      bool   `json:"success"`
	Location     string `json:"location"`
}

type CallbackResp struct {
	EBusinessID string `json:"e_business_id"`
	UpdateTime  string `json:"update_time"`
	Success     bool   `json:"success"`
	Reason      string `json:"reason"`
}

type QueryExpressResponse struct {
	EBusinessID  string    `json:"e_business_id"`
	ShipperCode  string    `json:"shipper_code"`
	Success      bool      `json:"success"`
	LogisticCode string    `json:"logistic_code"`
	State        string    `json:"state"`
	StateEx      string    `json:"state_ex"`
	Location     string    `json:"location"`
	Traces       []*Traces `json:"traces"`
	Raw          string    `json:"raw"`
}

type QueryExpressRequest struct {
	RequestData
	EBusinessID string `json:"e_business_id"`
	RequestType string `json:"request_type"`
	DataType    string `json:"data_type"`
	DataSign    string `json:"data_sign"`
}

type RequestData struct {
	ShipperCode  string `json:"shipper_code"`
	LogisticCode string `json:"logistic_code"`
	OrderCode    string `json:"order_code"`
}

type Request struct {
	Data        interface{} `json:"-"`           // 请求原始参数
	RequestData string      `json:"RequestData"` // 请求内容需进行URL(utf-8)编码。请求内容JSON格式，须和DataType一致。
	EBusinessID string      `json:"EBusinessID"` // 用户ID
	DataType    string      `json:"DataType"`    // 请求、返回数据类型：2-json(固定为2)；
	DataSign    string      `json:"DataSign"`    // 数据内容签名：把(请求内容(未编码)+ApiKey)进行MD5加密，然后Base64编码，最后进行URL(utf-8)编码
	RequestType string      `json:"RequestType"` // 请求指令类型：1002-轨迹即时查询
}

func (req *Request) ToUrlValues() url.Values {
	values := url.Values{}
	values.Set("RequestData", req.RequestData)
	values.Set("EBusinessID", req.EBusinessID)
	values.Set("DataType", req.DataType)
	values.Set("DataSign", req.DataSign)
	values.Set("RequestType", req.RequestType)
	return values
}

// =============================== 快递查询 ===================================

// ------------------------------- 查询请求 -----------------------------------

// QueryRequest 快递鸟查询单号请求参数 参考文档:https://www.yuque.com/kdnjishuzhichi/dfcrg1/yv7zgv
type QueryRequest struct {
	CustomerName *string      `json:"CustomerName,omitempty"` // 顺丰必填 需要寄件人/收件人的手机号后四位
	ShipperCode  *ShipperCode `json:"ShipperCode"`            // 快递公司编码
	LogisticCode *string      `json:"LogisticCode"`           // 快递单号
	Sort         *int         `json:"Sort,omitempty"`         // 排序方式：0-升序（默认），1-降序
	OrderCode    *string      `json:"OrderCode,omitempty"`    // 订单编号
}

// ------------------------------- 查询请求 -----------------------------------

// ------------------------------- 查询响应 -----------------------------------

// QueryResponse 快递鸟查询单号响应参数
type QueryResponse struct {
	EBusinessID    string      `json:"EBusinessID"`              // 用户ID
	ShipperCode    string      `json:"ShipperCode"`              // 快递公司编码
	LogisticCode   string      `json:"LogisticCode"`             // 快递单号
	Success        bool        `json:"Success"`                  // 成功与否 true/false
	Reason         string      `json:"Reason"`                   // 失败原因
	State          string      `json:"State"`                    // 物流状态 0-暂无轨迹信息,1-已揽收,2-在途中,3-签收,4-问题件,5-转寄,6-清关
	StateEx        DetailState `json:"StateEx"`                  // 物流状态的细分状态
	Location       string      `json:"Location"`                 // 快递当前城市
	Traces         []*Traces   `json:"Traces"`                   // 物流轨迹
	OrderCode      string      `json:"OrderCode,omitempty"`      // 订单编号
	Callback       string      `json:"Callback,omitempty"`       // 用户自定义回传字段
	Station        string      `json:"Station,omitempty"`        // 派件网点的名称
	StationTel     string      `json:"StationTel,omitempty"`     // 派件网点的电话
	StationAdd     string      `json:"StationAdd,omitempty"`     // 派件网点的地址
	DeliveryMan    string      `json:"DeliveryMan,omitempty"`    // 派件员姓名
	DeliveryManTel string      `json:"DeliveryManTel,omitempty"` // 派件员电话
	NextCity       string      `json:"NextCity,omitempty"`       // 下一站城市
}

// Traces 快递轨迹
type Traces struct {
	AcceptTime    string      `json:"AcceptTime"`    // 轨迹发生事件
	AcceptStation string      `json:"AcceptStation"` // 轨迹描述
	Location      string      `json:"Location"`      // 轨迹发生城市
	Action        DetailState `json:"Action"`        // 轨迹操作 通 StateEx 字段
	Remark        string      `json:"Remark"`        // 备注
}

// ------------------------------- 查询响应 -----------------------------------

// =============================== 快递查询 ===================================

// =============================== 轨迹订阅 ===================================

// ------------------------------- 订阅请求 -----------------------------------

// TrackSubscriptionRequest 轨迹订阅请求参数 参考文档: https://www.yuque.com/kdnjishuzhichi/dfcrg1/qkzowx
type TrackSubscriptionRequest struct {
	ShipperCode   ShipperCode `json:"ShipperCode"`             // 快递公司编码
	LogisticCode  string      `json:"LogisticCode"`            // 快递单号
	CustomerName  string      `json:"CustomerName,omitempty"`  // 顺丰必填 需要寄件人/收件人的手机号后四位
	Sort          int         `json:"Sort,omitempty"`          // 排序方式：0-升序（默认），1-降序
	OrderCode     string      `json:"OrderCode,omitempty"`     // 订单编号
	Callback      string      `json:"Callback,omitempty"`      // 用户自定义回传字段
	Receiver      *Receiver   `json:"Receiver,omitempty"`      // 收件人信息
	Sender        *Sender     `json:"Sender,omitempty"`        // 寄件人信息
	IsSendMessage bool        `json:"IsSendMessage,omitempty"` // 是否订阅短信通知 0-不需要，1-需要
}

// Receiver 收件人信息
type Receiver struct {
	Company      string `json:"Company"`      // 收件人公司
	Name         string `json:"Title"`        // 收件人
	Tel          string `json:"Tel"`          // 收件人电话
	Mobile       string `json:"Mobile"`       // 收件人手机 与电话号码二选一
	ProvinceName string `json:"ProvinceName"` // 收件人省份
	CityName     string `json:"CityName"`     // 收件人城市
	ExpAreaName  string `json:"ExpAreaName"`  // 收件人区域
	Address      string `json:"Address"`      // 收件人详细地址
}

// Sender 寄件人信息
type Sender struct {
	Company      string `json:"Company"`      // 寄件人公司
	Name         string `json:"Title"`        // 寄件人
	Tel          string `json:"Tel"`          // 寄件人电话
	Mobile       string `json:"Mobile"`       // 寄件人手机 与电话号码二选一
	ProvinceName string `json:"ProvinceName"` // 寄件人省份
	CityName     string `json:"CityName"`     // 寄件人城市
	ExpAreaName  string `json:"ExpAreaName"`  // 寄件人区域
	Address      string `json:"Address"`      // 寄件人详细地址
}

// TrackSubscriptionResponse 轨迹订阅响应参数
type TrackSubscriptionResponse struct {
	ShipperCode string `json:"ShipperCode"`
	UpdateTime  string `json:"UpdateTime"`
	EBusinessID string `json:"EBusinessID"`
	Success     bool   `json:"Success"`
}

// ------------------------------- 订阅请求 -----------------------------------

// ------------------------------- 订阅推送 -----------------------------------

// TrackSubPushRequest 轨迹订阅推送请求参数
type TrackSubPushRequest struct {
	RawRequestData string         `json:"-"`           // 原始数据
	RequestData    *PushTrackData `json:"RequestData"` // 请求参数
	DataSign       string         `json:"DataSign"`    // 数据签名
	RequestType    string         `json:"RequestType"` // 请求类型
}

// MapTrackSubPushRequest 轨迹订阅(地图)推送请求参数
type MapTrackSubPushRequest struct {
	RequestData *MapPushTrackData `json:"RequestData"` // 请求参数
	DataSign    string            `json:"DataSign"`    // 数据签名
	RequestType string            `json:"RequestType"` // 请求类型
}

// PushTrackData 轨迹推送数据
type PushTrackData struct {
	EBusinessID string       `json:"EBusinessID"` // 用户ID
	PushTime    string       `json:"PushTime"`    // 推送时间
	Data        []*TrackData `json:"Data"`        // 轨迹数据集合
	Count       string       `json:"Count"`       // 轨迹个数
}

// MapPushTrackData 轨迹(地图)推送数据
type MapPushTrackData struct {
	EBusinessID string          `json:"EBusinessID"` // 用户ID
	PushTime    string          `json:"PushTime"`    // 推送时间
	Data        []*MapTrackData `json:"Data"`        // 轨迹数据集合
	Count       string          `json:"Count"`       // 轨迹个数
}

// TrackData 轨迹数据集合
type TrackData struct {
	// 基础数据字段
	EBusinessID    string            `json:"EBusinessID"`              // 用户ID
	ShipperCode    string            `json:"ShipperCode"`              // 快递公司编码
	LogisticCode   string            `json:"LogisticCode"`             // 快递单号
	Success        bool              `json:"Success"`                  // 成功与否
	Reason         string            `json:"Reason,omitempty"`         // 失败原因
	OrderCode      string            `json:"OrderCode,omitempty"`      // 订单编号
	State          string            `json:"State"`                    // 物流状态
	StateEx        string            `json:"StateEx"`                  // 物流状态的细分状态
	Location       string            `json:"Location,omitempty"`       // 快递当前城市
	Callback       string            `json:"Callback,omitempty"`       // 用户自定义回传字段
	Traces         []*TrackDataTrace `json:"Traces"`                   // 物流轨迹
	Station        string            `json:"Station,omitempty"`        // 派件网点的名称
	StationTel     string            `json:"StationTel,omitempty"`     // 派件网点的电话
	StationAdd     string            `json:"StationAdd,omitempty"`     // 派件网点的地址
	DeliveryManTel string            `json:"DeliveryManTel,omitempty"` // 派件员电话
	DeliveryMan    string            `json:"DeliveryMan,omitempty"`    // 派件员姓名
	NextCity       string            `json:"NextCity,omitempty"`       // 下一站城市
}

type MapTrackData struct {
	TrackData
	// 轨迹地图推送数据字段
	ReceiverCityLatAndLng string          `json:"ReceiverCityLatAndLng"`    // 收件人城市经纬度
	SenderCityLatAndLng   string          `json:"SenderCityLatAndLng"`      // 寄件人城市经纬度
	Coordinates           *Coordinates    `json:"Coordinates,omitempty"`    // 当前城市经纬度
	PreCoordinates        *PreCoordinates `json:"PreCoordinates,omitempty"` // 预设路径经纬度
	RouteMapUrl           string          `json:"RouteMapUrl,omitempty"`    // 轨迹地图URL
}

// TrackDataTrace 轨迹数据
type TrackDataTrace struct {
	AcceptTime    string      `json:"AcceptTime"`         // 轨迹发生时间
	AcceptStation string      `json:"AcceptStation"`      // 轨迹描述
	Location      string      `json:"Location,omitempty"` // 轨迹发生城市
	Action        DetailState `json:"Action"`             // 轨迹操作
	Remark        string      `json:"Remark,omitempty"`   // 备注
}

// Coordinates 轨迹推送(地图)经纬度
type Coordinates struct {
	Location  string `json:"Location,omitempty"`  // 当前城市
	LatAndLng string `json:"LatAndLng,omitempty"` // 当前城市经纬度
}

// PreCoordinates 轨迹推送(地图)预设路径经纬度
type PreCoordinates struct {
	Location  string `json:"Location,omitempty"`  // 预设路径经过城市
	LatAndLng string `json:"LatAndLng,omitempty"` // 预设路径城市经纬度
}

// TrackSubPushResponse 轨迹订阅推送响应参数 https://www.yuque.com/kdnjishuzhichi/dfcrg1/meiubz#o8ytn
type TrackSubPushResponse struct {
	EBusinessID string `json:"EBusinessID"`      // 用户ID
	UpdateTime  string `json:"UpdateTime"`       // 更新时间
	Success     bool   `json:"Success"`          // 成功与否 true/false
	Reason      string `json:"Reason,omitempty"` // 失败原因
}

// ------------------------------- 订阅推送 -----------------------------------

// =============================== 轨迹订阅 ===================================
