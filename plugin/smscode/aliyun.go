package smscode

import (
	"context"
	"encoding/json"
	"fmt"

	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	"github.com/alibabacloud-go/dysmsapi-20170525/v2/client"
	"github.com/alibabacloud-go/tea/tea"

	"github.com/leor-w/kid/config"
	"github.com/leor-w/kid/plugin"
	"github.com/leor-w/kid/utils"
)

type Ali struct {
	Client  *client.Client
	options *AliOptions
}

type AliOption func(*AliOptions)

func (ali *Ali) Provide(ctx context.Context) interface{} {
	var configName string
	name, ok := ctx.Value(new(plugin.NameKey)).(string)
	if ok && len(name) > 0 {
		configName = "." + name
	}
	confPrefix := fmt.Sprintf("aliyun%s", configName)
	if !config.Exist(fmt.Sprintf(confPrefix)) {
		panic(fmt.Sprintf("config.yaml file not found configuration item [%s]", confPrefix))
	}
	return NewAliSMS(
		WithAccessKeyId(config.GetString(utils.GetConfigurationItem(confPrefix, "accessKeyId"))),
		WithAccessKeySecret(config.GetString(utils.GetConfigurationItem(confPrefix, "accessKeySecret"))),
		WithEndpoint(config.GetString(utils.GetConfigurationItem(confPrefix, "endpoint"))),
		WithSignName(config.GetString(utils.GetConfigurationItem(confPrefix, "signName"))),
		WithTempCode(config.GetString(utils.GetConfigurationItem(confPrefix, "tempCode"))),
	)
}

// SendSMS 发送短信验证码 默认参数名为 code
func (ali *Ali) SendSMS(phone, code string) error {
	return ali.SendSMSWithParams(phone, map[string]interface{}{"code": code})
}

// SendSMSWithParams 发送验证码指定参数名
func (ali *Ali) SendSMSWithParams(phone string, params map[string]interface{}) error {
	return ali.SendSMSWithTemplate(phone, ali.options.TempCode, params)
}

// SendSMSWithTemplate 发送验证码指定模版及参数
func (ali *Ali) SendSMSWithTemplate(phone string, template string, params map[string]interface{}) error {
	return ali.SendSMSWithSignAndTemplate(phone, ali.options.SignName, template, params)
}

// SendSMSWithSignAndTemplate 发送验证码指定短信签名、模版及参数 https://next.api.aliyun.com/api/Dysmsapi/2017-05-25/SendSms?lang=GO&params={}
func (ali *Ali) SendSMSWithSignAndTemplate(phone string, signName string, template string, params map[string]interface{}) error {
	paramsBytes, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("阿里云: 构建模版参数字段错误: %s", err.Error())
	}
	sendResp, err := ali.Client.SendSms(&client.SendSmsRequest{
		SignName:      &signName,
		TemplateCode:  &template,
		PhoneNumbers:  &phone,
		TemplateParam: tea.String(string(paramsBytes)),
	})
	if err != nil {
		return err
	}
	if *sendResp.Body.Code != "OK" {
		return fmt.Errorf("阿里云: 发送短信失败: %s", *sendResp.Body.Message)
	}
	return nil
}

func NewAliSMS(opts ...AliOption) SMSCode {
	var options AliOptions
	for _, o := range opts {
		o(&options)
	}
	ali := &Ali{
		options: &options,
	}
	conf := &openapi.Config{
		AccessKeyId:     tea.String(ali.options.AccessKeyId),
		AccessKeySecret: tea.String(ali.options.AccessKeySecret),
		Endpoint:        tea.String(ali.options.Endpoint),
	}
	cli, err := client.NewClient(conf)
	if err != nil {
		panic(fmt.Sprintf("阿里云: 初始化短信客户端失败: %s", err.Error()))
	}
	ali.Client = cli
	return ali
}
