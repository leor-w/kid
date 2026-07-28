package douyin

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"

	"github.com/lsjhtang/kid/logger"
	"github.com/lsjhtang/kid/utils"

	"github.com/qiniu/go-sdk/v7/sms/bytes"
)

type Payment struct {
	Options *PayOptions
}

const (
	ProductURL = "https://developer.toutiao.com/api/apps/ecpay/v1/"
	SandboxURL = "https://open-sandbox.douyin.com/api/apps/ecpay/v1/"
)

const (
	OtherSettleParams = "other_settle_params" // 其他分账方参数 (Other settle params)
	AppId             = "app_id"              // 小程序appID (Applets appID)
	ThirdpartyId      = "thirdparty_id"       // 代小程序进行该笔交易调用的第三方平台服务商 id (The id of the third-party platform service provider that calls the transaction on behalf of the Applets)
	Sign              = "sign"                // 签名 (sign)
)

var (
	CallbackCheckSignErr     = fmt.Errorf("回调请求验签失败") // 验签失败
	CallbackMessageDecodeErr = fmt.Errorf("回调请求信息解析错误")
)

// CreateOrder 预下单接口
// 文档: https://developer.open-douyin.com/docs/resource/zh-CN/mini-app/develop/server/ecpay/pay-list/pay
func (pay *Payment) CreateOrder(conf *PayCreateOrderConfig) (*PayCreateOrderResponse, error) {
	var resp PayCreateOrderResponse
	notifyUrl := utils.String(pay.Options.NotifyUrl)
	if conf.NotifyUrl != nil {
		notifyUrl = conf.NotifyUrl
	}
	req := &PayCreateOrderRequest{
		AppId:        utils.String(pay.Options.AppId),
		OutOrderNo:   conf.OutTradeNo,
		TotalAmount:  conf.TotalAmount,
		Subject:      conf.Subject,
		Body:         conf.Body,
		ValidTime:    conf.ValidTime,
		CpExtra:      conf.CpExtra,
		NotifyUrl:    notifyUrl,
		ThirdpartyId: utils.String(pay.Options.ThirdpartyId),
		StoreUid:     utils.String(pay.Options.StoreUid),
		DisableMsg:   conf.DisableMsg,
		MsgPage:      conf.MsgPage,
		LimitPayWay:  conf.LimitPayWay,
	}
	if conf.ExpandOrderInfo != nil {
		req.ExpandOrderInfo = conf.ExpandOrderInfo
	}
	req.Sign = utils.String(pay.sign(req.ToMap()))
	if err := pay.doPost("create_order", req, &resp); err != nil {
		return nil, err
	}
	if resp.ErrNo != 0 {
		return nil, fmt.Errorf("预下单失败: 错误码: [%d] 错误说明: [%s]", resp.ErrNo, resp.ErrTips)
	}
	return &resp, nil
}

func (pay *Payment) QueryOrder(outTradeNo string) (*QueryOrderResponse, error) {
	var resp QueryOrderResponse
	req := &QueryOrderRequest{
		AppId:      utils.String(pay.Options.AppId),
		OutOrderNo: &outTradeNo,
	}
	req.Sign = utils.String(pay.sign(req.ToMap()))
	if err := pay.doPost("query_order", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (pay *Payment) sign(params map[string]interface{}) string {
	var paramsArr []string
	for k, v := range params {
		if k == OtherSettleParams || k == AppId || k == ThirdpartyId || k == Sign {
			continue
		}
		value := strings.TrimSpace(fmt.Sprintf("%v", v))
		if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") && len(value) > 1 {
			value = value[1 : len(value)-1]
		}
		value = strings.TrimSpace(value)
		if value == "" || value == "null" || value == "<nil>" {
			continue
		}
		paramsArr = append(paramsArr, value)
	}

	paramsArr = append(paramsArr, pay.Options.Salt)
	sort.Strings(paramsArr)
	return fmt.Sprintf("%x", md5.Sum([]byte(strings.Join(paramsArr, "&"))))
}

func (pay *Payment) doPost(uri string, body, respData interface{}) error {
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return err
	}
	url := SandboxURL
	if pay.Options.IsProduct {
		url = ProductURL
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", url, uri), bytes.NewReader(bodyJson))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Errorf("close douyin response body error: %s", err.Error())
		}
	}(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http status code: %d", resp.StatusCode)
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(respBody, respData); err != nil {
		return err
	}
	return nil
}

// Callback 抖音回调请求处理
func (pay *Payment) Callback(req *CallbackReqeust) (*CallbackRequestMessage, error) {
	if ok := pay.callbackCheckSign(req); !ok {
		return nil, CallbackCheckSignErr
	}
	var msg CallbackRequestMessage
	if err := json.Unmarshal([]byte(req.Msg), &msg); err != nil {
		return nil, CallbackMessageDecodeErr
	}
	return &msg, nil
}

func (pay *Payment) callbackCheckSign(req *CallbackReqeust) bool {
	var reqArr = []string{req.Timestamp, req.Nonce, req.Msg, pay.Options.Token}
	sort.Strings(reqArr)
	h := sha1.New()
	h.Write([]byte(strings.Join(reqArr, "")))
	return req.MsgSignature == fmt.Sprintf("%x", h.Sum(nil))
}

func NewPayment(opts *PayOptions) *Payment {
	return &Payment{
		Options: opts,
	}
}
