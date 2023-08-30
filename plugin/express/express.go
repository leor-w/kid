package express

type Express interface {
	Query
	Consign
}

type Query interface {
	Query(params *QueryParams) (*QueryResponse, error)            // 实时查询
	QueryMapTrack(params *QueryParams) (*QueryMapResponse, error) // 实时查询地图轨迹

	Subscribe(params *SubscribeParams) error                              // 订阅轨迹
	VerifyCallbackPush(data, sign string) (*SubscribeCallbackPush, error) // 验证回调推送数据

	SubscribeMapTrack(params *SubscribeParams) error                            // 订阅地图轨迹
	VerifyCallbackMapPush(data, sign string) (*SubscribeCallbackMapPush, error) // 验证地图轨迹回调推送数据
}

type Consign interface {
	BShipmentOrder(params *SendExpressParams) (*SendExpressResponse, error) // B端寄件
	CShipmentOrder(params *SendExpressParams) (*SendExpressResponse, error) // C端寄件

	BShipmentCallbackVerify(sign, params string) (*BSendExpressCallbackPush, error) // B端寄件验证解析
	CShipmentCallbackVerify(sign, params string) (*CSendExpressCallbackPush, error) // C端寄件验证解析

	BCancelShipmentOrder(params *CancelExpressParams) (*CancelExpressResponse, error) // B端取消寄件
	CCancelShipmentOrder(params *CancelExpressParams) (*CancelExpressResponse, error) // C端取消寄件

	BShipmentExpressPushVerify(sign, params string) (*SubscribeCallbackPush, error) // B端寄件快递轨迹信息推送数据验证解析
	CShipmentExpressPushVerify(sign, params string) (*SubscribeCallbackPush, error) // C端寄件快递轨迹信息推送数据验证解析

	QueryBShipmentOrder(taskId string) (*QueryBSendExpressResponse, error) // 查询B端寄件信息

	QueryPostage(params *QueryPostageParams) (*QueryPostageResponse, error) // 查询运费
}
