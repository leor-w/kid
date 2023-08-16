package payments

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/consts"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/app"
)

type CombineApp struct {
	Client  *core.Client
	options *Options
}

func (payment *CombineApp) CloseOrder(ctx context.Context, req *CloseOrderRequest) (*core.APIResult, error) {
	var (
		localVarHTTPMethod   = http.MethodPost
		localVarPostBody     interface{}
		localVarQueryParams  url.Values
		localVarHeaderParams = http.Header{}
	)

	// Make sure Path Params are properly set
	if len(req.OutTradeNo) <= 0 {
		return nil, fmt.Errorf("field `OutTradeNo` is required and must be specified in CloseOrderRequest")
	}

	localVarPath := consts.WechatPayAPIServer + "/v3/combine-transactions/out-trade-no/{combine_out_trade_no}/close"
	// Build Path with Path Params
	localVarPath = strings.Replace(localVarPath, "{"+"out_trade_no"+"}",
		url.PathEscape(core.ParameterToString(req.OutTradeNo, "")), -1)

	// Make sure All Required Params are properly set

	// Setup Body Params
	localVarPostBody = &app.CloseRequest{
		Mchid: core.String(payment.options.Mchid),
	}

	// Determine the Content-Type Header
	localVarHTTPContentTypes := []string{"application/json"}
	// Setup Content-Type
	localVarHTTPContentType := core.SelectHeaderContentType(localVarHTTPContentTypes)

	return payment.Client.Request(ctx, localVarHTTPMethod, localVarPath, localVarHeaderParams,
		localVarQueryParams, localVarPostBody, localVarHTTPContentType)
}

func (payment *CombineApp) Prepay(ctx context.Context, req interface{}) (*Response, *core.APIResult, error) {
	var (
		localVarHTTPMethod   = http.MethodPost
		localVarPostBody     interface{}
		localVarQueryParams  url.Values
		localVarHeaderParams = http.Header{}
	)

	localVarPath := consts.WechatPayAPIServer + "/v3/combine-transactions/app"
	// Make sure All Required Params are properly set

	// Setup Body Params
	localVarPostBody = req

	// Determine the Content-Type Header
	localVarHTTPContentTypes := []string{"application/json"}
	// Setup Content-Type
	localVarHTTPContentType := core.SelectHeaderContentType(localVarHTTPContentTypes)

	// Perform Http Request
	result, err := payment.Client.Request(ctx, localVarHTTPMethod, localVarPath, localVarHeaderParams, localVarQueryParams, localVarPostBody, localVarHTTPContentType)
	if err != nil {
		return nil, result, err
	}

	// Extract PrepayResponse from Http Response
	resp := new(app.PrepayResponse)
	err = core.UnMarshalResponse(result.Response, resp)
	if err != nil {
		return nil, result, err
	}
	return &Response{Code: *resp.PrepayId}, result, nil
}

func (payment *CombineApp) QueryOrderById(ctx context.Context, req *QueryOrderByIdRequest) (*payments.Transaction, *core.APIResult, error) {
	return nil, nil, errors.New("not support method")
}

func (payment *CombineApp) QueryOrderByTradeNo(ctx context.Context, req *QueryOrderByOutTradeNoRequest) (*payments.Transaction, *core.APIResult, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		localVarQueryParams  url.Values
		localVarHeaderParams = http.Header{}
	)

	// Make sure Path Params are properly set
	if len(req.OutTradeNo) <= 0 {
		return nil, nil, fmt.Errorf("field `OutTradeNo` is required and must be specified in QueryOrderByOutTradeNoRequest")
	}

	localVarPath := consts.WechatPayAPIServer + "/v3/pay/transactions/out-trade-no/{out_trade_no}"
	// Build Path with Path Params
	localVarPath = strings.Replace(localVarPath, "{"+"out_trade_no"+"}",
		url.PathEscape(core.ParameterToString(req.OutTradeNo, "")), -1)

	// Make sure All Required Params are properly set
	if len(payment.options.Mchid) <= 0 {
		return nil, nil, fmt.Errorf("field `Mchid` is required and must be specified in QueryOrderByOutTradeNoRequest")
	}

	// Setup Query Params
	localVarQueryParams = url.Values{}
	localVarQueryParams.Add("mchid", core.ParameterToString(payment.options.Mchid, ""))

	// Determine the Content-Type Header
	localVarHTTPContentTypes := []string{}
	// Setup Content-Type
	localVarHTTPContentType := core.SelectHeaderContentType(localVarHTTPContentTypes)

	// Perform Http Request
	result, err := payment.Client.Request(ctx, localVarHTTPMethod, localVarPath, localVarHeaderParams,
		localVarQueryParams, localVarPostBody, localVarHTTPContentType)
	if err != nil {
		return nil, result, err
	}

	// Extract payments.Transaction from Http Response
	resp := new(payments.Transaction)
	err = core.UnMarshalResponse(result.Response, resp)
	if err != nil {
		return nil, result, err
	}
	return resp, result, nil
}
