package smscode

import (
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	"github.com/alibabacloud-go/dysmsapi-20170525/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/leor-w/kid/config"
	"github.com/leor-w/kid/logger"
)

type Ali struct {
	Client  *client.Client
	options *AliOptions
}

func (ali *Ali) Provide() interface{} {
	return NewAliSMS(
		WithAccessKeyId(config.GetString("aliyun.accessKeyId")),
		WithAccessKeySecret(config.GetString("aliyun.accessKeySecret")),
		WithEndpoint(config.GetString("aliyun.endpoint")),
		WithSignName(config.GetString("aliyun.signName")),
		WithTempCode(config.GetString("aliyun.tempCode")),
	)
}

func (ali *Ali) Init() error {
	conf := &openapi.Config{
		AccessKeyId:     tea.String(ali.options.AccessKeyId),
		AccessKeySecret: tea.String(ali.options.AccessKeySecret),
		Endpoint:        tea.String(ali.options.Endpoint),
	}
	cli, err := client.NewClient(conf)
	if err != nil {
		return err
	}
	ali.Client = cli
	return nil
}

// SendSMS 发送验证码 https://next.api.aliyun.com/api/Dysmsapi/2017-05-25/SendSms?lang=GO&params={}
func (ali *Ali) SendSMS(phone, text string) error {
	sendResp, err := ali.Client.SendSms(&client.SendSmsRequest{
		SignName:      tea.String(ali.options.SignName),
		TemplateCode:  tea.String(ali.options.TempCode),
		PhoneNumbers:  &phone,
		TemplateParam: tea.String(fmt.Sprintf("{\"code\": \"%s\"}", text)),
	})
	if err != nil {
		return err
	}
	if *sendResp.Body.Code != "OK" {
		return fmt.Errorf("aliyun send validate code sms failed: %s", *sendResp.Body.Message)
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
	if err := ali.Init(); err != nil {
		logger.Errorf("init aliyun sms service failed: %s", err.Error())
	}
	return ali
}

type AliOption func(*AliOptions)

type AliOptions struct {
	AccessKeyId     string
	AccessKeySecret string
	Endpoint        string
	SignName        string
	TempCode        string
}

func WithAccessKeyId(accessKeyId string) AliOption {
	return func(o *AliOptions) {
		o.AccessKeyId = accessKeyId
	}
}

func WithAccessKeySecret(accessKeySecret string) AliOption {
	return func(o *AliOptions) {
		o.AccessKeySecret = accessKeySecret
	}
}

func WithEndpoint(endpoint string) AliOption {
	return func(o *AliOptions) {
		o.Endpoint = endpoint
	}
}

func WithSignName(signName string) AliOption {
	return func(o *AliOptions) {
		o.SignName = signName
	}
}

func WithTempCode(tempCode string) AliOption {
	return func(o *AliOptions) {
		o.TempCode = tempCode
	}
}
