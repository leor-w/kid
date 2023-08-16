package payments

import (
	"encoding/json"
	"fmt"
)

type CombineAppPrepayRequest struct {
	// 合单商户appid
	CombineAppid *string `json:"combine_appid"`
	// 合单商户号
	CombineMchid *string `json:"combine_mchid"`
	// 合单商户订单号
	CombineOutTradeNo *string `json:"combine_out_trade_no"`
	// 场景信息
	SceneInfo []SceneInfo `json:"scene_info,omitempty"`
	// 子单信息
	SubOrders []SubOrder `json:"sub_orders"`
	// 支付者
	CombinePayerInfo []PayerInfo `json:"combine_payer_info,omitempty"`
	// 交易起始时间
	TimeStart *string `json:"time_start,omitempty"`
	// 交易结束时间
	TimeExpire *string `json:"time_expire,omitempty"`
	// 通知地址
	NotifyUrl *string `json:"notify_url"`
}

func (c CombineAppPrepayRequest) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if c.CombineAppid == nil {
		return nil, fmt.Errorf("field `CombineAppid` is required and must be specified in CombineAppPrepayRequest")
	}
	toSerialize["combine_appid"] = c.CombineAppid

	if c.CombineMchid == nil {
		return nil, fmt.Errorf("field `CombineMchid` is required and must be specified in CombineAppPrepayRequest")
	}
	toSerialize["combine_mchid"] = c.CombineMchid

	if c.CombineOutTradeNo == nil {
		return nil, fmt.Errorf("field `CombineOutTradeNo` is required and must be specified in CombineAppPrepayRequest")
	}
	toSerialize["combine_out_trade_no"] = c.CombineOutTradeNo

	if c.SceneInfo != nil {
		toSerialize["scene_info"] = c.SceneInfo
	}

	if c.SubOrders == nil {
		return nil, fmt.Errorf("field `SubOrders` is required and must be specified in CombineAppPrepayRequest")
	}
	toSerialize["sub_orders"] = c.SubOrders

	if c.CombinePayerInfo != nil {
		toSerialize["combine_payer_info"] = c.CombinePayerInfo
	}

	if c.TimeStart != nil {
		toSerialize["time_start"] = c.TimeStart
	}

	if c.TimeExpire != nil {
		toSerialize["time_expire"] = c.TimeExpire
	}

	if c.NotifyUrl == nil {
		return nil, fmt.Errorf("field `NotifyUrl` is required and must be specified in CombineAppPrepayRequest")
	}
	toSerialize["notify_url"] = c.NotifyUrl
	return json.Marshal(toSerialize)
}

func (c CombineAppPrepayRequest) String() string {
	var ret string
	if c.CombineAppid == nil {
		ret += "CombineAppid:<nil>, "
	} else {
		ret += fmt.Sprintf("CombineAppid:%v, ", c.CombineAppid)
	}

	if c.CombineMchid == nil {
		ret += "CombineMchid:<nil>, "
	} else {
		ret += fmt.Sprintf("CombineMchid:%v, ", c.CombineMchid)
	}

	if c.CombineOutTradeNo == nil {
		ret += "CombineOutTradeNo:<nil>, "
	} else {
		ret += fmt.Sprintf("CombineOutTradeNo:%v, ", c.CombineOutTradeNo)
	}

	if c.SceneInfo == nil {
		ret += "SceneInfo:<nil>, "
	} else {
		ret += fmt.Sprintf("SceneInfo:%v, ", c.SceneInfo)
	}

	if c.SubOrders == nil {
		ret += "SubOrders:<nil>, "
	} else {
		ret += fmt.Sprintf("SubOrders:%v, ", c.SubOrders)
	}

	if c.CombinePayerInfo == nil {
		ret += "CombinePayerInfo:<nil>, "
	} else {
		ret += fmt.Sprintf("CombinePayerInfo:%v, ", c.CombinePayerInfo)
	}

	if c.TimeStart == nil {
		ret += "TimeStart:<nil>, "
	} else {
		ret += fmt.Sprintf("TimeStart:%v, ", c.TimeStart)
	}

	if c.TimeExpire == nil {
		ret += "TimeExpire:<nil>, "
	} else {
		ret += fmt.Sprintf("TimeExpire:%v, ", c.TimeExpire)
	}

	if c.NotifyUrl == nil {
		ret += "NotifyUrl:<nil>, "
	} else {
		ret += fmt.Sprintf("NotifyUrl:%v, ", c.NotifyUrl)
	}
	return fmt.Sprintf("CombineAppPrepayRequest:{%s}", ret)
}

func (c CombineAppPrepayRequest) Clone() *CombineAppPrepayRequest {
	var ret CombineAppPrepayRequest
	if c.CombineAppid != nil {
		ret.CombineAppid = new(string)
		*ret.CombineAppid = *c.CombineAppid
	}
	if c.CombineMchid != nil {
		ret.CombineMchid = new(string)
		*ret.CombineMchid = *c.CombineMchid
	}
	if c.CombineOutTradeNo != nil {
		ret.CombineOutTradeNo = new(string)
		*ret.CombineOutTradeNo = *c.CombineOutTradeNo
	}
	if c.SceneInfo != nil {
		ret.SceneInfo = make([]SceneInfo, len(c.SceneInfo))
		for i, item := range c.SceneInfo {
			ret.SceneInfo[i] = *item.Clone()
		}
	}
	if c.SubOrders != nil {
		ret.SubOrders = make([]SubOrder, len(c.SubOrders))
		for i, item := range c.SubOrders {
			ret.SubOrders[i] = *item.Clone()
		}
	}
	if c.CombinePayerInfo != nil {
		ret.CombinePayerInfo = make([]PayerInfo, len(c.CombinePayerInfo))
		for i, item := range c.CombinePayerInfo {
			ret.CombinePayerInfo[i] = *item.Clone()
		}
	}
	if c.TimeStart != nil {
		ret.TimeStart = new(string)
		*ret.TimeStart = *c.TimeStart
	}
	if c.TimeExpire != nil {
		ret.TimeExpire = new(string)
		*ret.TimeExpire = *c.TimeExpire
	}
	if c.NotifyUrl != nil {
		ret.NotifyUrl = new(string)
		*ret.NotifyUrl = *c.NotifyUrl
	}
	return &ret
}

// SubOrder 子单信息
type SubOrder struct {
	// 子单商户号
	Mchid *string `json:"mchid"`
	// 附加数据
	Attach *string `json:"attach"`
	// 订单金额
	Amount *Amount `json:"amount"`
	// 子单商户订单号
	OutTradeNo *string `json:"out_trade_no"`
	// 订单优惠标记
	GoodsTag *string `json:"goods_tag,omitempty"`
	// 商品描述
	Description *string     `json:"description"`
	SettleInfo  *SettleInfo `json:"settle_info,omitempty"`
}

func (s SubOrder) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if s.Mchid == nil {
		return nil, fmt.Errorf("field `Mchid` is required and must be specified in SubOrder")
	}
	toSerialize["mchid"] = s.Mchid

	if s.Attach == nil {
		return nil, fmt.Errorf("field `Attach` is required and must be specified in SubOrder")
	}
	toSerialize["attach"] = s.Attach

	if s.Amount == nil {
		return nil, fmt.Errorf("field `Amount` is required and must be specified in SubOrder")
	}
	toSerialize["amount"] = s.Amount

	if s.OutTradeNo == nil {
		return nil, fmt.Errorf("field `OutTradeNo` is required and must be specified in SubOrder")
	}
	toSerialize["out_trade_no"] = s.OutTradeNo
	if s.GoodsTag != nil {
		toSerialize["goods_tag"] = s.GoodsTag
	}
	if s.Description == nil {
		return nil, fmt.Errorf("field `Description` is required and must be specified in SubOrder")
	}
	toSerialize["description"] = s.Description
	if s.SettleInfo != nil {
		toSerialize["settle_info"] = s.SettleInfo
	}
	return json.Marshal(toSerialize)
}

func (s SubOrder) String() string {
	var ret string
	if s.Mchid == nil {
		ret += "Mchid:<nil>, "
	} else {
		ret += fmt.Sprintf("Mchid:%v, ", s.Mchid)
	}
	if s.Attach == nil {
		ret += "Attach:<nil>, "
	} else {
		ret += fmt.Sprintf("Attach:%v, ", s.Attach)
	}
	if s.Amount == nil {
		ret += "Amount:<nil>, "
	} else {
		ret += fmt.Sprintf("Amount:%v, ", s.Amount)
	}
	if s.OutTradeNo == nil {
		ret += "OutTradeNo:<nil>, "
	} else {
		ret += fmt.Sprintf("OutTradeNo:%v, ", s.OutTradeNo)
	}
	if s.GoodsTag == nil {
		ret += "GoodsTag:<nil>, "
	} else {
		ret += fmt.Sprintf("GoodsTag:%v, ", s.GoodsTag)
	}
	if s.Description == nil {
		ret += "Description:<nil>, "
	} else {
		ret += fmt.Sprintf("Description:%v, ", s.Description)
	}
	if s.SettleInfo == nil {
		ret += "SettleInfo:<nil>, "
	} else {
		ret += fmt.Sprintf("SettleInfo:%v, ", s.SettleInfo)
	}
	return fmt.Sprintf("SubOrder:{%s}", ret)
}

func (s SubOrder) Clone() *SubOrder {
	var ret SubOrder
	if s.Mchid != nil {
		ret.Mchid = new(string)
		*ret.Mchid = *s.Mchid
	}
	if s.Attach != nil {
		ret.Attach = new(string)
		*ret.Attach = *s.Attach
	}
	if s.Amount != nil {
		ret.Amount = new(Amount)
		*ret.Amount = *s.Amount.Clone()
	}
	if s.OutTradeNo != nil {
		ret.OutTradeNo = new(string)
		*ret.OutTradeNo = *s.OutTradeNo
	}
	if s.GoodsTag != nil {
		ret.GoodsTag = new(string)
		*ret.GoodsTag = *s.GoodsTag
	}
	if s.Description != nil {
		ret.Description = new(string)
		*ret.Description = *s.Description
	}
	if s.SettleInfo != nil {
		ret.SettleInfo = new(SettleInfo)
		*ret.SettleInfo = *s.SettleInfo.Clone()
	}
	return nil
}

// SettleInfo
type SettleInfo struct {
	// 是否指定分账
	ProfitSharing *bool `json:"profit_sharing,omitempty"`
	// 补差金额
	SubsidyAmount *int64 `json:"subsidy_amount,omitempty"`
}

func (o SettleInfo) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}

	if o.ProfitSharing != nil {
		toSerialize["profit_sharing"] = o.ProfitSharing
	}
	if o.SubsidyAmount != nil {
		toSerialize["subsidy_amount"] = o.SubsidyAmount
	}
	return json.Marshal(toSerialize)
}

func (o SettleInfo) String() string {
	var ret string
	if o.ProfitSharing == nil {
		ret += "ProfitSharing:<nil>"
	} else {
		ret += fmt.Sprintf("ProfitSharing:%v", *o.ProfitSharing)
	}
	if o.SubsidyAmount == nil {
		ret += "SubsidyAmount:<nil>"
	} else {
		ret += fmt.Sprintf("SubsidyAmount:%v", *o.SubsidyAmount)
	}

	return fmt.Sprintf("SettleInfo{%s}", ret)
}

func (o SettleInfo) Clone() *SettleInfo {
	ret := SettleInfo{}

	if o.ProfitSharing != nil {
		ret.ProfitSharing = new(bool)
		*ret.ProfitSharing = *o.ProfitSharing
	}

	if o.SubsidyAmount != nil {
		ret.SubsidyAmount = new(int64)
		*ret.SubsidyAmount = *o.SubsidyAmount
	}

	return &ret
}

type Amount struct {
	// 订单总金额，单位为分
	Total *int64 `json:"total"`
	// CNY：人民币，境内商户号仅支持人民币。
	Currency *string `json:"currency,omitempty"`
}

func (o Amount) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}

	if o.Total == nil {
		return nil, fmt.Errorf("field `Total` is required and must be specified in Amount")
	}
	toSerialize["total"] = o.Total

	if o.Currency != nil {
		toSerialize["currency"] = o.Currency
	}
	return json.Marshal(toSerialize)
}

func (o Amount) String() string {
	var ret string
	if o.Total == nil {
		ret += "Total:<nil>, "
	} else {
		ret += fmt.Sprintf("Total:%v, ", *o.Total)
	}

	if o.Currency == nil {
		ret += "Currency:<nil>"
	} else {
		ret += fmt.Sprintf("Currency:%v", *o.Currency)
	}

	return fmt.Sprintf("Amount{%s}", ret)
}

func (o Amount) Clone() *Amount {
	ret := Amount{}

	if o.Total != nil {
		ret.Total = new(int64)
		*ret.Total = *o.Total
	}

	if o.Currency != nil {
		ret.Currency = new(string)
		*ret.Currency = *o.Currency
	}

	return &ret
}

type PayerInfo struct {
	OpenId *string `json:"open_id,omitempty"`
}

func (p PayerInfo) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if p.OpenId != nil {
		toSerialize["open_id"] = p.OpenId
	}
	return json.Marshal(toSerialize)
}

func (p PayerInfo) String() string {
	var ret string
	if p.OpenId == nil {
		ret += "OpenId:<nil>, "
	} else {
		ret += fmt.Sprintf("OpenId:%v, ", *p.OpenId)
	}
	return fmt.Sprintf("PayerInfo: {%s}", ret)
}

func (p PayerInfo) Clone() *PayerInfo {
	ret := PayerInfo{}

	if p.OpenId != nil {
		ret.OpenId = new(string)
		*ret.OpenId = *p.OpenId
	}

	return &ret
}

// SceneInfo 支付场景描述
type SceneInfo struct {
	// 用户终端IP
	PayerClientIp *string `json:"payer_client_ip"`
	// 商户端设备号
	DeviceId  *string    `json:"device_id,omitempty"`
	StoreInfo *StoreInfo `json:"store_info,omitempty"`
}

func (o SceneInfo) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}

	if o.PayerClientIp == nil {
		return nil, fmt.Errorf("field `PayerClientIp` is required and must be specified in SceneInfo")
	}
	toSerialize["payer_client_ip"] = o.PayerClientIp

	if o.DeviceId != nil {
		toSerialize["device_id"] = o.DeviceId
	}

	if o.StoreInfo != nil {
		toSerialize["store_info"] = o.StoreInfo
	}
	return json.Marshal(toSerialize)
}

func (o SceneInfo) String() string {
	var ret string
	if o.PayerClientIp == nil {
		ret += "PayerClientIp:<nil>, "
	} else {
		ret += fmt.Sprintf("PayerClientIp:%v, ", *o.PayerClientIp)
	}

	if o.DeviceId == nil {
		ret += "DeviceId:<nil>, "
	} else {
		ret += fmt.Sprintf("DeviceId:%v, ", *o.DeviceId)
	}

	return fmt.Sprintf("SceneInfo{%s}", ret)
}

func (o SceneInfo) Clone() *SceneInfo {
	ret := SceneInfo{}

	if o.PayerClientIp != nil {
		ret.PayerClientIp = new(string)
		*ret.PayerClientIp = *o.PayerClientIp
	}

	if o.DeviceId != nil {
		ret.DeviceId = new(string)
		*ret.DeviceId = *o.DeviceId
	}

	if o.StoreInfo != nil {
		ret.StoreInfo = o.StoreInfo.Clone()
	}

	return &ret
}

// StoreInfo 商户门店信息
type StoreInfo struct {
	// 商户侧门店编号
	Id *string `json:"id"`
	// 商户侧门店名称
	Name *string `json:"name,omitempty"`
	// 地区编码，详细请见微信支付提供的文档
	AreaCode *string `json:"area_code,omitempty"`
	// 详细的商户门店地址
	Address *string `json:"address,omitempty"`
}

func (o StoreInfo) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}

	if o.Id == nil {
		return nil, fmt.Errorf("field `Id` is required and must be specified in StoreInfo")
	}
	toSerialize["id"] = o.Id

	if o.Name != nil {
		toSerialize["name"] = o.Name
	}

	if o.AreaCode != nil {
		toSerialize["area_code"] = o.AreaCode
	}

	if o.Address != nil {
		toSerialize["address"] = o.Address
	}
	return json.Marshal(toSerialize)
}

func (o StoreInfo) String() string {
	var ret string
	if o.Id == nil {
		ret += "Id:<nil>, "
	} else {
		ret += fmt.Sprintf("Id:%v, ", *o.Id)
	}

	if o.Name == nil {
		ret += "Name:<nil>, "
	} else {
		ret += fmt.Sprintf("Name:%v, ", *o.Name)
	}

	if o.AreaCode == nil {
		ret += "AreaCode:<nil>, "
	} else {
		ret += fmt.Sprintf("AreaCode:%v, ", *o.AreaCode)
	}

	if o.Address == nil {
		ret += "Address:<nil>"
	} else {
		ret += fmt.Sprintf("Address:%v", *o.Address)
	}

	return fmt.Sprintf("StoreInfo{%s}", ret)
}

func (o StoreInfo) Clone() *StoreInfo {
	ret := StoreInfo{}

	if o.Id != nil {
		ret.Id = new(string)
		*ret.Id = *o.Id
	}

	if o.Name != nil {
		ret.Name = new(string)
		*ret.Name = *o.Name
	}

	if o.AreaCode != nil {
		ret.AreaCode = new(string)
		*ret.AreaCode = *o.AreaCode
	}

	if o.Address != nil {
		ret.Address = new(string)
		*ret.Address = *o.Address
	}

	return &ret
}
