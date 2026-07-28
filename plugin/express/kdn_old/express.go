package express

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/lsjhtang/kid/injector"

	"github.com/lsjhtang/kid/config"
	"github.com/lsjhtang/kid/logger"
	"github.com/lsjhtang/kid/utils"
)

// Express 快递100 SDK
type Express struct {
	options *Options
}

func (exp *Express) Provide(ctx context.Context) interface{} {
	var confName string
	name, ok := ctx.Value(injector.NameKey{}).(string)
	if ok && len(name) > 0 {
		confName = "." + name
	}
	confPrefix := fmt.Sprintf("kd100_old%s", confName)
	if !config.Exist(confPrefix) {
		panic(fmt.Sprintf("config.yaml file not found configuration item [%s]", confPrefix))
	}
	return New(
		WithKey(config.GetString(utils.GetConfigurationItem(confPrefix, "key"))),
		WithCustomer(config.GetString(utils.GetConfigurationItem(confPrefix, "customer"))),
		WithSecret(config.GetString(utils.GetConfigurationItem(confPrefix, "secret"))),
		WithBaseUrl(config.GetString(utils.GetConfigurationItem(confPrefix, "baseUrl"))),
		WithSalt(config.GetString(utils.GetConfigurationItem(confPrefix, "salt"))),
		WithNotifyUrl(config.GetString(utils.GetConfigurationItem(confPrefix, "notifyUrl"))),
	)
}

type Option func(*Options)

// QueryExpress 快递查询 https://api.kuaidi100.com/document/5f0ffb5ebc8da837cbd8aefc
func (exp *Express) QueryExpress(conf *QueryReqConfig) (*QueryExpressResp, error) {
	var queryReq = &QueryParams{
		Com:      conf.Com,
		Num:      conf.Num,
		Phone:    conf.Phone,
		From:     conf.From,
		To:       conf.To,
		Resultv2: conf.Resultv2,
		Show:     conf.Show,
		Order:    conf.Order,
	}
	params, err := json.Marshal(queryReq)
	if err != nil {
		return nil, fmt.Errorf("序列化请求参数失败：%s", err.Error())
	}
	values := url.Values{}
	values.Set("customer", exp.options.Customer)
	values.Set("sign", exp.sign(string(params)+exp.options.Key+exp.options.Customer))
	values.Set("param", string(params))
	var resp QueryExpressResp
	if err := exp.doPost("/poll/query.do", &values, &resp); err != nil {
		return nil, fmt.Errorf("查询失败：%s", err.Error())
	}
	if resp.ReturnCode != "200" {
		return nil, fmt.Errorf("查询失败: 错误码[%s], 错误信息[%s]", resp.ReturnCode, resp.Message)
	}
	return &resp, nil
}

// Subscription 订阅快递 https://api.kuaidi100.com/document/5f0ffa8f2977d50a94e1023c
func (exp *Express) Subscription(conf *SubscribeReqConfig) error {
	callbackUrl := conf.CallbackUrl
	if len(callbackUrl) <= 0 {
		callbackUrl = exp.options.NotifyUrl
	}
	var subscribeReq = &SubscribeParam{
		Company: conf.Company,
		Number:  conf.Number,
		From:    conf.From,
		To:      conf.To,
		Key:     exp.options.Key,
		Parameters: &SubscribeParameters{
			CallbackUrl:        callbackUrl,
			Salt:               exp.options.Salt,
			Resultv2:           conf.Resultv2,
			AutoCom:            conf.AutoCom,
			InterCom:           conf.InterCom,
			DepartureCountry:   conf.DepartureCountry,
			DepartureCom:       conf.DepartureCom,
			DestinationCountry: conf.DestinationCountry,
			DestinationCom:     conf.DestinationCom,
			Phone:              conf.Phone,
		},
	}
	params, err := json.Marshal(subscribeReq)
	if err != nil {
		return fmt.Errorf("序列化请求参数失败：%s", err.Error())
	}
	values := url.Values{}
	values.Set("schema", conf.Schema)
	values.Set("param", string(params))
	values.Set("sign", exp.sign(string(params)+exp.options.Salt))
	var subscribeResp SubscribeResp
	if err := exp.doPost("/poll", &values, &subscribeResp); err != nil {
		return fmt.Errorf("订阅失败：%s", err.Error())
	}
	if subscribeResp.Result != true {
		return fmt.Errorf("订阅失败: 错误码: [%s], 错误信息: [%s]", subscribeResp.ReturnCode, subscribeResp.Message)
	}
	return nil
}

// AnalyticalPushData 解析推送数据 https://api.kuaidi100.com/document/5f0ffa8f2977d50a94e1023c
func (exp *Express) AnalyticalPushData(data, sign string) (*SubscribePush, error) {
	if err := exp.verify(data, sign); err != nil {
		return nil, fmt.Errorf("验证签名失败：%s", err.Error())
	}
	var push SubscribePush
	if err := json.Unmarshal([]byte(data), &push); err != nil {
		return nil, fmt.Errorf("解析推送数据失败：%s", err.Error())
	}
	return &push, nil
}

// doPost 发送 POST 请求
func (exp *Express) doPost(uri string, values *url.Values, respData interface{}) error {
	resp, err := http.Post(exp.options.BaseUrl+uri, "application/x-www-form-urlencoded", strings.NewReader(values.Encode()))
	if err != nil {
		return fmt.Errorf("请求失败：%s", err.Error())
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Errorf("关闭响应失败：%s", err.Error())
		}
	}(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("请求失败，状态码：%d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败：%s", err.Error())
	}
	if respData != nil {
		if err := json.Unmarshal(body, respData); err != nil {
			return fmt.Errorf("解析响应失败：%s", err.Error())
		}
	}
	return nil
}

// sign 生成签名
func (exp *Express) sign(waitSign string) string {
	hash := md5.Sum([]byte(waitSign))
	return strings.ToUpper(fmt.Sprintf("%x", hash))
}

// verify 验证签名
func (exp *Express) verify(data, sign string) error {
	if exp.sign(data+exp.options.Salt) != sign {
		return fmt.Errorf("签名错误")
	}
	return nil
}

// IsRequiredPhone 是否必须需要手机号
func IsRequiredPhone(expCom string) bool {
	coms := []string{SF, SFKY, FWSY}
	for _, com := range coms {
		if com == expCom {
			return true
		}
	}
	return false
}

func New(opts ...Option) *Express {
	var options = &Options{}
	for _, opt := range opts {
		opt(options)
	}
	return &Express{options: options}
}
