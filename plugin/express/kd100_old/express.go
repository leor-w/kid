package kdn

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/lsjhtang/kid/injector"

	"github.com/lsjhtang/kid/config"
	"github.com/lsjhtang/kid/logger"
	"github.com/lsjhtang/kid/utils"
)

type IExpress interface {
	QueryByShipperCode(req *QueryRequest) (*QueryResponse, error)                        // 即时查询
	Query(response *QueryResponse) (*QueryResponse, error)                               // 快递查询
	TrackSubscription(req *TrackSubscriptionRequest) (*TrackSubscriptionResponse, error) // 轨迹订阅
	TrackSubscriptionResponse(status bool, reason string) *TrackSubPushResponse          // 轨迹订阅推送请求响应
	AnalyticalTrackPushData(data string) (*TrackSubPushRequest, error)                   // 解析轨迹推送数据
	Verify(reqData, sign string) error                                                   // 验证推送数据
}

type Express struct {
	Options *Options
}

func (express *Express) Provide(ctx context.Context) interface{} {
	var confName string
	name, ok := ctx.Value(injector.NameKey{}).(string)
	if ok && len(name) > 0 {
		confName = "." + name
	}
	confPrefix := fmt.Sprintf("kdn%s", confName)
	if !config.Exist(confPrefix) {
		panic(fmt.Sprintf("config.yaml file not found configuration item [%s]", confPrefix))
	}
	return New(
		WithEBusinessID(config.GetString(utils.GetConfigurationItem(confPrefix, "eBusinessID"))),
		WithAppKey(config.GetString(utils.GetConfigurationItem(confPrefix, "appKey"))),
		WithDataType(config.GetString(utils.GetConfigurationItem(confPrefix, "dataType"))),
		WithRequestType(config.GetString(utils.GetConfigurationItem(confPrefix, "requestType"))),
		WithBaseURL(config.GetString(utils.GetConfigurationItem(confPrefix, "baseURL"))),
	)
}

type Option func(o *Options)

func New(opts ...Option) *Express {
	option := &Options{}
	for _, opt := range opts {
		opt(option)
	}
	express := &Express{
		Options: option,
	}
	return express
}

// QueryByShipperCode 即时查询(增值版) https://www.yuque.com/kdnjishuzhichi/dfcrg1/yv7zgv#JRGGr
// 通过快递公司编码查询，需要传 ShipperCode，LogisticCode。并且按单计费。
func (express *Express) QueryByShipperCode(req *QueryRequest) (*QueryResponse, error) {
	var resp QueryResponse
	if err := express.doPost("/Ebusiness/EbusinessOrderHandle.aspx", &Request{
		Data:        req,
		EBusinessID: express.Options.EBusinessID,
		DataType:    "2",
		RequestType: "8001",
	}, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Query 即时查询(增值版) https://www.yuque.com/kdnjishuzhichi/dfcrg1/yv7zgv#JRGGr
// 通过运单号查询，不需要传 ShipperCode，但是需要传 LogisticCode，即运单号。并且按次数计费。
func (express *Express) Query(response *QueryResponse) (*QueryResponse, error) {
	var req = &Request{
		Data:        response,
		EBusinessID: express.Options.EBusinessID,
		DataType:    "2",
		RequestType: "8002",
	}
	var resp QueryResponse
	if err := express.doPost("/Ebusiness/EbusinessOrderHandle.aspx", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// TrackSubscription 轨迹订阅(增值版) https://www.yuque.com/kdnjishuzhichi/dfcrg1/qkzowx#TYqFA
func (express *Express) TrackSubscription(req *TrackSubscriptionRequest) (*TrackSubscriptionResponse, error) {
	var resp TrackSubscriptionResponse
	if err := express.doPost("/api/dist", &Request{
		Data:        req,
		EBusinessID: express.Options.EBusinessID,
		DataType:    "2",
		RequestType: "8008",
	}, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// TrackSubscriptionResponse 轨迹订阅推送请求响应 https://www.yuque.com/kdnjishuzhichi/dfcrg1/qkzowx#TYqFA
func (express *Express) TrackSubscriptionResponse(status bool, reason string) *TrackSubPushResponse {
	return &TrackSubPushResponse{
		EBusinessID: express.Options.EBusinessID,
		UpdateTime:  time.Now().Format("2006-01-02 15:04:05"),
		Success:     status,
		Reason:      reason,
	}
}

// AnalyticalTrackPushData 解析轨迹推送数据 https://www.yuque.com/kdnjishuzhichi/dfcrg1/tx5q8r#JRGGr
func (express *Express) AnalyticalTrackPushData(data string) (*TrackSubPushRequest, error) {
	var trackPushData TrackSubPushRequest
	if err := json.Unmarshal([]byte(data), &trackPushData); err != nil {
		return nil, fmt.Errorf("解析推送数据出错：%w", err)
	}
	return &trackPushData, nil
}

// doPost 发送 POST 请求 https://www.yuque.com/kdnjishuzhichi/dfcrg1/zes04h
func (express *Express) doPost(uri string, req *Request, respData interface{}) error {
	reqData, err := json.Marshal(req.Data)
	if err != nil {
		return fmt.Errorf("序列化请求数据出错：%w", err)
	}
	req.RequestData = string(reqData)
	req.DataSign = url.QueryEscape(express.sign(req.RequestData))
	resp, err := http.PostForm(express.Options.BaseURL+uri, req.ToUrlValues())
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Errorf("关闭响应出错：%s", err.Error())
		}
	}(resp.Body)
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应出错：%w", err)
	}
	if err := json.Unmarshal(data, respData); err != nil {
		return fmt.Errorf("解析响应出错：%w", err)
	}
	return nil
}

// sign 生成数据签名 https://www.yuque.com/kdnjishuzhichi/dfcrg1/zes04h
// 原始请求数据+APIkey进行MD5加密，然后Base64编码，最后进行URL（utf-8）编码
func (express *Express) sign(reqData string) string {
	h := md5.New()
	io.WriteString(h, reqData+express.Options.AppKey)
	//h.Write([]byte(reqData + express.Options.AppKey))
	cipherStr := fmt.Sprintf("%x", h.Sum(nil))
	base64Str := base64Encode([]byte(cipherStr))
	return string(base64Str)
}

func base64Encode(src []byte) []byte {
	return []byte(base64.StdEncoding.EncodeToString(src))
}

// Verify 验证回调数据签名 https://www.yuque.com/kdnjishuzhichi/dfcrg1/zes04h
func (express *Express) Verify(reqData, sign string) error {
	waitSign := reqData + express.Options.AppKey
	signed := express.sign(waitSign)
	if signed != sign {
		return fmt.Errorf("签名验证失败")
	}
	return nil
}
