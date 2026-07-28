package kd100

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/spf13/cast"

	"github.com/lsjhtang/kid/plugin/express"
)

type Express struct {
	cli *Client

	options *Options
}

type Option func(options *Options)

// Query 快递查询 https://api.kuaidi100.com/document/5f0ffb5ebc8da837cbd8aefc
func (exp *Express) Query(params *express.QueryParams) (*express.QueryResponse, error) {
	var queryReq = &QueryParams{
		Com:      string(GetStandardMapCompanyCode(params.CompCode)),
		Num:      params.ExpNo,
		Phone:    params.Phone,
		From:     params.From,
		To:       params.To,
		Resultv2: cast.ToString(params.CustomParams[ResultV2]),
		Show:     "0",
		Order:    params.Sort,
	}
	paramsJson, err := json.Marshal(queryReq)
	if err != nil {
		return nil, fmt.Errorf("序列化请求参数失败：%s", err.Error())
	}
	paramsStr := string(paramsJson)
	var values url.Values
	values.Set("customer", exp.options.Customer)
	values.Set("sign", exp.Sign(paramsStr+exp.options.Key+exp.options.Customer))
	values.Set("param", paramsStr)

	var resp QueryExpressResp
	if err := exp.cli.DoPost("/poll/query.do", &values, &resp); err != nil {
		return nil, fmt.Errorf("查询失败：%s", err.Error())
	}
	if resp.ReturnCode != "200" {
		return nil, fmt.Errorf("查询失败: 错误码[%s], 错误信息[%s]", resp.ReturnCode, resp.Message)
	}

	var tracks []*express.Track
	for _, datum := range resp.Data {
		trackTime, err := time.Parse("2006-01-02 15:04:05", datum.Time)
		if err != nil {
			return nil, fmt.Errorf("解析时间 [%s] 失败：%s", datum.Time, err.Error())
		}
		tracks = append(tracks, &express.Track{
			Description: datum.Context,
			Time:        trackTime,
			AreaCode:    datum.AreaCode,
			Area:        datum.AreaName,
			Status:      GetStatusToStandardState(datum.Status),
			Location:    datum.Location,
			Coordinates: datum.AreaCenter,
		})
	}
	return &express.QueryResponse{
		State:    GetStatusToStandardState(resp.State),
		IsSigned: cast.ToBool(resp.Ischeck),
		CompCode: GetStandardCode(resp.Com),
		ExpNo:    resp.Nu,
		Tracks:   tracks,
	}, nil
}

// QueryMapTrack 查询物流轨迹地图 https://api.kuaidi100.com/document/5ff2c3e7ba1bf00302f5612e
func (exp *Express) QueryMapTrack(params *express.QueryParams) (*express.QueryMapResponse, error) {
	query := &QueryMapTrackParams{
		QueryParams: QueryParams{
			Com:      string(GetStandardMapCompanyCode(params.CompCode)),
			Num:      params.ExpNo,
			Phone:    params.Phone,
			From:     params.From,
			To:       params.To,
			Resultv2: cast.ToString(params.CustomParams[ResultV2]),
			Show:     "0",
			Order:    params.Sort,
		},
		MapConfigKey: exp.options.MapTempKey,
	}
	paramsJson, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("序列化请求参数失败：%s", err.Error())
	}
	paramsStr := string(paramsJson)
	var values url.Values
	values.Set("customer", exp.options.Customer)
	values.Set("sign", exp.Sign(paramsStr+exp.options.Key+exp.options.Customer))
	values.Set("param", paramsStr)
	var resp QueryMapTrackResponse
	if err := exp.cli.DoPost("/poll/maptrack.do", &values, &resp); err != nil {
		return nil, fmt.Errorf("查询失败：%s", err.Error())
	}
	if resp.Status != "200" {
		return nil, fmt.Errorf("查询失败: 错误码[%s], 错误信息[%s]", resp.Status, resp.Message)
	}
	arrivalTime, err := time.Parse("2006-01-02 15:04:05", resp.ArrivalTime)
	if err != nil {
		return nil, fmt.Errorf("解析时间 [%s] 失败：%s", resp.ArrivalTime, err.Error())
	}
	var tracks []*express.QueryMapResponseData
	for _, datum := range resp.Data {
		trackTime, err := time.Parse("2006-01-02 15:04:05", datum.Time)
		if err != nil {
			return nil, fmt.Errorf("解析时间 [%s] 失败：%s", datum.Time, err.Error())
		}
		tracks = append(tracks, &express.QueryMapResponseData{
			Describe:   datum.Context,
			Time:       trackTime,
			AreaCode:   datum.AreaCode,
			Area:       datum.AreaName,
			State:      GetStatusToStandardState(resp.State),
			Coordinate: datum.AreaCenter,
			Location:   datum.Location,
		})
	}

	return &express.QueryMapResponse{
		State:       GetStatusToStandardState(resp.State),
		IsSigned:    cast.ToBool(resp.IsCheck),
		CompCode:    GetStandardCode(resp.Com),
		ExpNo:       resp.Nu,
		TrailUrl:    resp.TrailUrl,
		ArrivalTime: arrivalTime,
		TotalTime:   resp.TotalTime,
		RemainTime:  resp.RemainTime,
		Data:        tracks,
	}, nil
}

// Subscribe 订阅快递 https://api.kuaidi100.com/document/5f0ffa8f2977d50a94e1023c
func (exp *Express) Subscribe(params *express.SubscribeParams) error {
	callbackUrl := params.CallbackUrl
	if len(callbackUrl) <= 0 {
		callbackUrl = exp.options.NotifyUrl
	}
	var subscribeReq = &SubscribeParam{
		Company: string(GetStandardMapCompanyCode(params.CompCode)),
		Number:  params.ExpNo,
		From:    params.From,
		To:      params.To,
		Key:     exp.options.Key,
		Parameters: &SubscribeParameters{
			CallbackUrl: callbackUrl,
			Salt:        exp.options.Salt,
			Resultv2:    cast.ToString(params.CustomParams[ResultV2]),
			Phone:       params.Phone,
		},
	}
	paramsJson, err := json.Marshal(subscribeReq)
	if err != nil {
		return fmt.Errorf("序列化请求参数失败：%s", err.Error())
	}
	values := url.Values{}
	values.Set("schema", params.Schema)
	values.Set("param", string(paramsJson))
	var subscribeResp SubscribeResp
	if err := exp.cli.DoPost("/poll", &values, &subscribeResp); err != nil {
		return fmt.Errorf("订阅失败：%s", err.Error())
	}
	if subscribeResp.Result != true {
		return fmt.Errorf("订阅失败: 错误码: [%s], 错误信息: [%s]", subscribeResp.ReturnCode, subscribeResp.Message)
	}
	return nil
}

// VerifyCallbackPush 验证回调推送 https://api.kuaidi100.com/document/5f0ffa8f2977d50a94e1023c#section_1
func (exp *Express) VerifyCallbackPush(data, sign string) (*express.SubscribeCallbackPush, error) {
	if exp.Verify(sign, data) {
		return nil, fmt.Errorf("验证签名失败")
	}
	var push SubscribePush
	if err := json.Unmarshal([]byte(data), &push); err != nil {
		return nil, fmt.Errorf("解析推送数据失败：%s", err.Error())
	}
	isSigned := false
	if push.LastResult.Ischeck == "1" {
		isSigned = true
	}
	var tracks []*express.Track
	for _, datum := range push.LastResult.Data {
		trackTime, err := time.Parse("2006-01-02 15:04:05", datum.Time)
		if err != nil {
			return nil, fmt.Errorf("解析时间 [%s] 失败：%s", datum.Time, err.Error())
		}
		tracks = append(tracks, &express.Track{
			Description: datum.Context,
			Time:        trackTime,
			AreaCode:    datum.AreaCode,
			Area:        datum.AreaName,
			Status:      GetStatusToStandardState(datum.StatusCode),
			Location:    datum.Location,
			Coordinates: datum.AreaCenter,
		})
	}
	return &express.SubscribeCallbackPush{
		Status:            push.Status,
		Message:           push.Message,
		ShouldResubscribe: false,
		State:             GetStatusToStandardState(push.LastResult.State),
		IsSigned:          isSigned,
		CompCode:          GetStandardCode(push.LastResult.Com),
		ExpNo:             push.LastResult.Nu,
		Tracks:            tracks,
	}, nil
}

// SubscribeMapTrack 订阅物流轨迹地图 https://api.kuaidi100.com/document/603f47dfa62a19500e19866f#section_0
func (exp *Express) SubscribeMapTrack(params *express.SubscribeParams) error {
	callbackUrl := params.CallbackUrl
	if len(callbackUrl) <= 0 {
		callbackUrl = exp.options.NotifyUrl
	}
	var subscribeReq = &SubscribeMapTrackParam{
		Company: string(GetStandardMapCompanyCode(params.CompCode)),
		Number:  params.ExpNo,
		From:    params.From,
		To:      params.To,
		Key:     exp.options.Key,
		Parameters: &SubscribeMapTrackParameters{
			CallbackUrl:  callbackUrl,
			Salt:         exp.options.Salt,
			Phone:        params.Phone,
			MapConfigKey: exp.options.MapTempKey,
			ResultV2:     cast.ToString(params.CustomParams[ResultV2]),
		},
	}
	paramsJson, err := json.Marshal(subscribeReq)
	if err != nil {
		return fmt.Errorf("序列化请求参数失败：%s", err.Error())
	}
	values := url.Values{}
	values.Set("schema", params.Schema)
	values.Set("param", string(paramsJson))
	var subscribeResp SubscribeResp
	if err := exp.cli.DoPost("/poll", &values, &subscribeResp); err != nil {
		return fmt.Errorf("订阅失败：%s", err.Error())
	}
	if subscribeResp.Result != true {
		return fmt.Errorf("订阅失败: 错误码: [%s], 错误信息: [%s]", subscribeResp.ReturnCode, subscribeResp.Message)
	}
	return nil
}

// VerifyCallbackMapPush 验证回调推送 https://api.kuaidi100.com/document/603f47dfa62a19500e19866f#section_1
func (exp *Express) VerifyCallbackMapPush(data, sign string) (*express.SubscribeCallbackMapPush, error) {
	if exp.Verify(sign, data) {
		return nil, fmt.Errorf("验证签名失败")
	}
	var push SubscribeMapTrackPush
	if err := json.Unmarshal([]byte(data), &push); err != nil {
		return nil, fmt.Errorf("解析推送数据失败：%s", err.Error())
	}
	isSigned := false
	if push.LastResult.Ischeck == "1" {
		isSigned = true
	}
	var tracks []*express.Track
	for _, datum := range push.LastResult.Data {
		trackTime, err := time.Parse("2006-01-02 15:04:05", datum.Time)
		if err != nil {
			return nil, fmt.Errorf("解析时间 [%s] 失败：%s", datum.Time, err.Error())
		}
		tracks = append(tracks, &express.Track{
			Description: datum.Context,
			Time:        trackTime,
			AreaCode:    datum.AreaCode,
			Area:        datum.AreaName,
			Status:      GetStatusToStandardState(datum.StatusCode),
			Location:    datum.Location,
			Coordinates: datum.AreaCenter,
		})
	}
	arrivalTime, err := time.Parse("2006-01-02 15:04:05", push.LastResult.ArrivalTime)
	if err != nil {
		return nil, fmt.Errorf("解析预计到达时间 [%s] 失败：%s", push.LastResult.ArrivalTime, err.Error())
	}
	return &express.SubscribeCallbackMapPush{
		Status:            push.Status,
		Message:           push.Message,
		ShouldResubscribe: false,
		State:             GetStatusToStandardState(push.LastResult.State),
		IsSigned:          isSigned,
		CompCode:          GetStandardCode(push.LastResult.Com),
		ExpNo:             push.LastResult.Nu,
		TrailUrl:          push.LastResult.TrailUrl,
		ArrivalTime:       arrivalTime,
		TotalTime:         push.LastResult.TotalTime,
		RemainTime:        push.LastResult.RemainTime,
		Tracks:            tracks,
	}, nil
}

// BShipmentOrder 商家寄快递 https://api.kuaidi100.com/document/603cb649a62a19500e19866b#section_0
func (exp *Express) BShipmentOrder(params *express.SendExpressParams) (*express.SendExpressResponse, error) {
	var callback int
	if params.Callback {
		callback = 1
	}
	var req = &BCreateShipmentOrderReq{
		Kuaidicom:        string(GetStandardMapCompanyCode(params.CompCode)),
		RecManName:       params.Recipient,
		RecManMobile:     params.RecipientPhone,
		RecManPrintAddr:  params.RecipientAddr,
		SendManName:      params.SenderName,
		SendManMobile:    params.SenderPhone,
		SendManPrintAddr: params.SenderAddr,
		Cargo:            params.ItemName,
		CallBackUrl:      params.CallbackUrl,
		Payment:          "SHIPPER",
		ServiceType:      string(GetCompanyCodeServiceType(GetStandardMapCompanyCode(params.CompCode))),
		Weight:           cast.ToString(params.Weight),
		Remark:           params.Remark,
		Salt:             exp.options.Salt,
		DayType:          params.DayType,
		ValinsPay:        cast.ToString(params.InsuredAmount),
		Op:               cast.ToString(callback),
		PollCallBackUrl:  exp.options.SendExpressNotifyUrl,
		Resultv2:         cast.ToString(params.CustomParams[ResultV2]),
	}
	paramsJson, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("序列化请求参数失败：%s", err.Error())
	}
	paramsStr := string(paramsJson)
	values := url.Values{}
	values.Set("method", "bOrder")
	values.Set("key", exp.options.Key)
	values.Set("sign", exp.Sign(paramsStr+exp.options.Key+exp.options.Secret))
	values.Set("t", cast.ToString(time.Now().UnixMilli()))
	values.Set("param", paramsStr)
	var resp BCreateShipmentOrderResp
	if err := exp.cli.DoPost("/order/borderapi.do", &values, &resp); err != nil {
		return nil, fmt.Errorf("创建订单失败：%s", err.Error())
	}
	if resp.Result != true {
		return nil, fmt.Errorf("创建订单失败: 错误码: [%s], 错误信息: [%s]", resp.ReturnCode, resp.Message)
	}
	return &express.SendExpressResponse{
		TaskID:  resp.Data.TaskId,
		OrderNo: resp.Data.OrderId,
		ExpNo:   resp.Data.Kuaidinum,
	}, nil
}

// CShipmentOrder 客户寄快递 https://api.kuaidi100.com/document/cduan-ji-jian-jie-kou-wen-dang#section_0
func (exp *Express) CShipmentOrder(params *express.SendExpressParams) (*express.SendExpressResponse, error) {
	callback := "0"
	if params.Callback {
		callback = "1"
	}
	var req = &CShipmentOrderReq{
		Kuaidicom:        string(GetStandardMapCompanyCode(params.CompCode)),
		RecManName:       params.Recipient,
		RecManMobile:     params.RecipientPhone,
		RecManPrintAddr:  params.RecipientAddr,
		SendManName:      params.SenderName,
		SendManMobile:    params.SenderPhone,
		SendManPrintAddr: params.SenderAddr,
		CallBackUrl:      params.CallbackUrl,
		Cargo:            params.ItemName,
		Payment:          GetPayType(params.PayType),
		Weight:           cast.ToString(params.Weight),
		Remark:           params.Remark,
		DayType:          params.DayType,
		PickupStartTime:  params.StartTime,
		PickupEndTime:    params.EndTime,
		Salt:             exp.options.Salt,
		Op:               callback,
		PollCallBackUrl:  params.CallbackUrl,
		Resultv2:         cast.ToString(params.CustomParams[ResultV2]),
		ThirdOrderId:     params.TradeNo,
	}
	paramsJson, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("序列化请求参数失败：%s", err.Error())
	}
	paramsStr := string(paramsJson)
	values := url.Values{}
	values.Set("method", "cOrder")
	values.Set("key", exp.options.Key)
	values.Set("sign", exp.Sign(paramsStr+exp.options.Key+exp.options.Secret))
	values.Set("t", cast.ToString(time.Now().UnixMilli()))
	values.Set("param", paramsStr)
	var resp CShipmentOrderResp
	if err := exp.cli.DoPost("/order/corderapi.do", &values, &resp); err != nil {
		return nil, fmt.Errorf("创建订单失败：%s", err.Error())
	}
	if resp.Result != true {
		return nil, fmt.Errorf("创建订单失败: 错误码: [%s], 错误信息: [%s]", resp.ReturnCode, resp.Message)
	}
	return &express.SendExpressResponse{
		TaskID:  resp.Data.TaskId,
		OrderNo: resp.Data.OrderId,
	}, nil
}

// BShipmentCallbackVerify 商家下单回调接口 https://api.kuaidi100.com/document/603cb649a62a19500e19866b#section_1
func (exp *Express) BShipmentCallbackVerify(sign, params string) (*express.BSendExpressCallbackPush, error) {
	if !exp.Verify(sign, params) {
		return nil, fmt.Errorf("验证签名失败")
	}
	var push CreateShipmentOrderPush
	if err := json.Unmarshal([]byte(params), &push); err != nil {
		return nil, fmt.Errorf("解析推送数据失败：%s", err.Error())
	}
	var fees []*express.FeeDetail
	for _, detail := range push.Data.FeeDetails {
		fees = append(fees, &express.FeeDetail{
			FeeType:   detail.FeeType,
			FeeDesc:   detail.FeeDesc,
			Fee:       detail.Amount,
			PayStatus: detail.PayStatus,
		})
	}
	return &express.BSendExpressCallbackPush{
		CompCode:         GetStandardCode(push.Kuaidicom),
		ExpNo:            push.Kuaidinum,
		OrderId:          push.Data.OrderId,
		Status:           push.Data.Status,
		PickupAgentName:  push.Data.CourierName,
		PickupAgentPhone: push.Data.CourierMobile,
		Weight:           cast.ToFloat64(push.Data.Weight),
		Freight:          cast.ToFloat64(push.Data.Freight),
		Volume:           cast.ToFloat64(push.Data.Volume),
		ActualWeight:     cast.ToFloat64(push.Data.ActualWeight),
		Fees:             fees,
	}, nil
}

// CShipmentCallbackVerify 客户下单回调接口 https://api.kuaidi100.com/document/cduan-ji-jian-jie-kou-wen-dang#section_1
func (exp *Express) CShipmentCallbackVerify(sign, params string) (*express.CSendExpressCallbackPush, error) {
	if exp.Verify(sign, params) {
		return nil, fmt.Errorf("验证签名失败")
	}
	var push CShipmentOrderCallback
	if err := json.Unmarshal([]byte(params), &push); err != nil {
		return nil, fmt.Errorf("解析推送数据失败：%s", err.Error())
	}

	return &express.CSendExpressCallbackPush{
		CompCode:         GetStandardCode(push.Kuaidicom),
		ExpNo:            push.Kuaidinum,
		OrderId:          push.Data.OrderId,
		Status:           cast.ToInt(push.Status),
		PickupAgentName:  push.Data.CourierName,
		PickupAgentPhone: push.Data.CourierMobile,
		Weight:           cast.ToFloat64(push.Data.Weight),
		Freight:          cast.ToFloat64(push.Data.Freight),
	}, nil
}

// BCancelShipmentOrder 商家寄件下单取消接口 https://api.kuaidi100.com/document/603cb649a62a19500e19866b#section_2
func (exp *Express) BCancelShipmentOrder(params *express.CancelExpressParams) (*express.CancelExpressResponse, error) {
	var req = &BCancelShipmentOrderReq{
		OrderId:   params.OrderId,
		TaskId:    params.TaskId,
		CancelMsg: params.Reason,
	}
	paramsJson, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("序列化请求参数失败：%s", err.Error())
	}
	paramsStr := string(paramsJson)
	values := url.Values{}
	values.Set("method", "cancel")
	values.Set("key", exp.options.Key)
	values.Set("sign", exp.Sign(paramsStr+exp.options.Key+exp.options.Secret))
	values.Set("t", cast.ToString(time.Now().UnixMilli()))
	values.Set("param", paramsStr)
	var resp BCancelShipmentOrderResp
	if err := exp.cli.DoPost("/order/borderapi.do", &values, &resp); err != nil {
		return nil, fmt.Errorf("取消订单失败：%s", err.Error())
	}
	if resp.Result != true {
		return nil, fmt.Errorf("取消订单失败: 错误码: [%s], 错误信息: [%s]", resp.ReturnCode, resp.Message)
	}
	return &express.CancelExpressResponse{
		Success: resp.Result,
		Code:    resp.ReturnCode,
		Msg:     resp.Message,
	}, nil
}

// CCancelShipmentOrder 客户寄件下单取消接口 https://api.kuaidi100.com/document/cduan-ji-jian-jie-kou-wen-dang#section_2
func (exp *Express) CCancelShipmentOrder(params *express.CancelExpressParams) (*express.CancelExpressResponse, error) {
	var req = &CShipmentOrderCancelReq{
		TaskId:    params.TaskId,
		OrderId:   params.OrderId,
		CancelMsg: params.Reason,
	}
	paramsJson, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("序列化请求参数失败：%s", err.Error())
	}
	paramsStr := string(paramsJson)
	values := url.Values{}
	values.Set("method", "cancel")
	values.Set("key", exp.options.Key)
	values.Set("sign", exp.Sign(paramsStr+exp.options.Key+exp.options.Secret))
	values.Set("t", cast.ToString(time.Now().UnixMilli()))
	values.Set("param", paramsStr)
	var resp CShipmentOrderCancelResp
	if err := exp.cli.DoPost("/order/corderapi.do", &values, &resp); err != nil {
		return nil, fmt.Errorf("取消订单失败：%s", err.Error())
	}
	if !resp.Result {
		return nil, fmt.Errorf("取消订单失败: 错误码: [%s], 错误信息: [%s]", resp.ReturnCode, resp.Message)
	}
	return &express.CancelExpressResponse{
		Success: resp.Result,
		Code:    resp.ReturnCode,
		Msg:     resp.Message,
	}, nil
}

// BShipmentExpressPushVerify 快递信息推送接口 https://api.kuaidi100.com/document/603cb649a62a19500e19866b#section_4
func (exp *Express) BShipmentExpressPushVerify(sign, params string) (*express.SubscribeCallbackPush, error) {
	if !exp.Verify(sign, params) {
		return nil, fmt.Errorf("验证签名失败")
	}
	var resp BSendExpressOrderPush
	if err := json.Unmarshal([]byte(params), &resp); err != nil {
		return nil, fmt.Errorf("解析推送数据失败：%s", err.Error())
	}
	var tracks []*express.Track
	for _, datum := range resp.LastResult.Data {
		trackTime, err := time.Parse("2006-01-02 15:04:05", datum.Time)
		if err != nil {
			return nil, fmt.Errorf("解析时间 [%s] 失败：%s", datum.Time, err.Error())
		}
		tracks = append(tracks, &express.Track{
			Description: datum.Context,
			Time:        trackTime,
			AreaCode:    datum.AreaCode,
			Area:        datum.AreaName,
			Status:      GetStatusToStandardState(datum.StatusCode),
			Location:    datum.Location,
			Coordinates: datum.AreaCenter,
		})
	}

	var isSigned bool
	if resp.LastResult.Ischeck == "1" {
		isSigned = true
	}

	return &express.SubscribeCallbackPush{
		Status:            resp.Status,
		Message:           resp.Message,
		ShouldResubscribe: false,
		State:             GetStatusToStandardState(resp.LastResult.State),
		IsSigned:          isSigned,
		CompCode:          GetStandardCode(resp.LastResult.Com),
		ExpNo:             resp.LastResult.Nu,
		Tracks:            tracks,
	}, nil
}

func (exp *Express) CShipmentExpressPushVerify(sign, params string) (*express.SubscribeCallbackPush, error) {
	if exp.Verify(sign, params) {
		return nil, fmt.Errorf("验证签名失败")
	}
	var push CExpressCallbackPush
	if err := json.Unmarshal([]byte(params), &push); err != nil {
		return nil, fmt.Errorf("解析推送数据失败：%s", err.Error())
	}
	var isSigned bool
	if push.LastResult.Ischeck == "1" {
		isSigned = true
	}
	var tracks []*express.Track
	for _, datum := range push.LastResult.Data {
		trackTime, err := time.Parse("2006-01-02 15:04:05", datum.Time)
		if err != nil {
			return nil, fmt.Errorf("解析时间 [%s] 失败：%s", datum.Time, err.Error())
		}
		tracks = append(tracks, &express.Track{
			Description: datum.Context,
			Time:        trackTime,
			AreaCode:    datum.AreaCode,
			Area:        datum.AreaName,
			Status:      GetStatusToStandardState(datum.StatusCode),
			Location:    datum.Location,
			Coordinates: datum.AreaCenter,
		})
	}
	return &express.SubscribeCallbackPush{
		Status:            push.Status,
		Message:           push.Message,
		ShouldResubscribe: false,
		State:             GetStatusToStandardState(push.LastResult.State),
		IsSigned:          isSigned,
		CompCode:          GetStandardCode(push.LastResult.Com),
		ExpNo:             push.LastResult.Nu,
		Tracks:            tracks,
	}, nil
}

// QueryBShipmentOrder 商家寄件信息查询接口 https://api.kuaidi100.com/document/603cb649a62a19500e19866b#section_8
func (exp *Express) QueryBShipmentOrder(taskId string) (*express.QueryBSendExpressResponse, error) {
	var req = &BQueryShipmentOrderReq{
		TaskId: taskId,
	}
	paramsJson, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("序列化请求参数失败：%s", err.Error())
	}
	paramsStr := string(paramsJson)
	values := url.Values{}
	values.Set("method", "query")
	values.Set("key", exp.options.Key)
	values.Set("sign", exp.Sign(paramsStr+exp.options.Key+exp.options.Secret))
	values.Set("t", cast.ToString(time.Now().UnixMilli()))
	values.Set("param", paramsStr)
	var resp BQueryShipmentOrderResp
	if err := exp.cli.DoPost("/order/borderapi.do", &values, &resp); err != nil {
		return nil, fmt.Errorf("查询订单失败：%s", err.Error())
	}
	if resp.Result != true {
		return nil, fmt.Errorf("查询订单失败: 错误码: [%s], 错误信息: [%s]", resp.ReturnCode, resp.Message)
	}

	createdTime, err := time.Parse("2006-01-02 15:04:05", resp.Data.CreateTime)
	if err != nil {
		return nil, fmt.Errorf("解析创建时间 [%s] 失败：%s", resp.Data.CreateTime, err.Error())
	}

	var details []*express.FeeDetail
	for _, detail := range resp.Data.FeeDetails {
		details = append(details, &express.FeeDetail{
			FeeType:   detail.FeeType,
			FeeDesc:   detail.FeeDesc,
			Fee:       detail.Amount,
			PayStatus: detail.PayStatus,
		})
	}
	return &express.QueryBSendExpressResponse{
		TaskId:               resp.Data.TaskId,
		CreatedAt:            createdTime,
		OrderId:              resp.Data.OrderId,
		CompCode:             GetStandardCode(resp.Data.KuaidiCom),
		ExpNo:                resp.Data.KuaidiNum,
		SenderName:           resp.Data.SendName,
		SenderPhone:          resp.Data.SendMobile,
		SenderProvince:       resp.Data.SendProvince,
		SenderCity:           resp.Data.SendCity,
		SenderDistrict:       resp.Data.SendDistrict,
		SenderAddress:        resp.Data.SendAddr,
		ReceiverName:         resp.Data.RecName,
		ReceiverPhone:        resp.Data.RecMobile,
		ReceiverProvince:     resp.Data.RecProvince,
		ReceiverCity:         resp.Data.RecCity,
		ReceiverDistrict:     resp.Data.RecDistrict,
		ReceiverAddress:      resp.Data.RecAddr,
		Valins:               resp.Data.Valins,
		ItemName:             resp.Data.Cargo,
		ExpressType:          resp.Data.ServiceType,
		Weight:               cast.ToFloat64(resp.Data.PreWeight),
		ActualWeight:         cast.ToFloat64(resp.Data.LastWeight),
		Remark:               resp.Data.Comment,
		ReservationData:      resp.Data.DayType,
		ReservationStartTime: resp.Data.PickupStartTime,
		ReservationEndTime:   resp.Data.PickupEndTime,
		PayType:              GetPayTypeReverse(resp.Data.Payment),
		Freight:              cast.ToFloat64(resp.Data.Freight),
		PickupAgentName:      resp.Data.CourierName,
		PickupAgentPhone:     resp.Data.CourierMobile,
		Status:               resp.Data.Status,
		PayStatus:            cast.ToUint8(resp.Data.PayStatus),
		FeeDetails:           details,
	}, nil
}

// QueryPostage 查询运费 https://api.kuaidi100.com/document/603cb649a62a19500e19866b#section_3
func (exp *Express) QueryPostage(params *express.QueryPostageParams) (*express.QueryPostageResponse, error) {
	var req = &QueryShipmentPriceReq{
		Kuaidicom:        string(GetStandardMapCompanyCode(params.CompCode)),
		SendManPrintAddr: params.SenderAddr,
		RecManPrintAddr:  params.ReceiverAddr,
		Weight:           cast.ToString(params.Weight),
		ServiceType:      string(GetCompanyCodeServiceType(GetStandardMapCompanyCode(params.CompCode))),
	}
	paramsJson, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("序列化请求参数失败：%s", err.Error())
	}
	paramsStr := string(paramsJson)
	values := url.Values{}
	values.Set("method", "price")
	values.Set("key", exp.options.Key)
	values.Set("sign", exp.Sign(paramsStr+exp.options.Key+exp.options.Secret))
	values.Set("t", cast.ToString(time.Now().UnixMilli()))
	values.Set("param", paramsStr)
	var resp QueryShipmentPriceResp
	if err := exp.cli.DoPost("/order/borderapi.do", &values, &resp); err != nil {
		return nil, fmt.Errorf("查询运费失败：%s", err.Error())
	}
	if resp.Result != true {
		return nil, fmt.Errorf("查询运费失败: 错误码: [%s], 错误信息: [%s]", resp.ReturnCode, resp.Message)
	}
	return &express.QueryPostageResponse{
		FirstWeight:     cast.ToFloat64(resp.Data.FirstPrice),
		ContinuedWeight: cast.ToFloat64(resp.Data.OverPrice),
		TotalPrice:      cast.ToFloat64(resp.Data.Price),
		ExpType:         resp.Data.ServiceType,
	}, nil
}
