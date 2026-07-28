package alicloud

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/lsjhtang/kid/database/redis"
	"github.com/lsjhtang/kid/logger"

	"github.com/lsjhtang/kid/plugin/lock"

	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	"github.com/alibabacloud-go/dysmsapi-20170525/v2/client"
	"github.com/alibabacloud-go/tea/tea"

	"github.com/lsjhtang/kid/config"
	"github.com/lsjhtang/kid/plugin"
	"github.com/lsjhtang/kid/plugin/sms"
	"github.com/lsjhtang/kid/utils"
)

type Adapter struct {
	Client  *client.Client
	db      *redis.Client
	lock    lock.Lock
	options *AliOptions
}

type AliOption func(*AliOptions)

func (ali *Adapter) Provide(ctx context.Context) any {
	var configName string
	name, ok := ctx.Value(new(plugin.NameKey)).(string)
	if ok && len(name) > 0 {
		configName = "." + name
	}
	confPrefix := fmt.Sprintf("aliyun%s", configName)
	if !config.Exist(fmt.Sprintf(confPrefix)) {
		panic(fmt.Sprintf("config.yaml file not found configuration item [%s]", confPrefix))
	}
	return NewAliCloudAdapter(
		WithAnalog(config.GetBool(utils.GetConfigurationItem(confPrefix, "analog"))),
		WithAccessKeyId(config.GetString(utils.GetConfigurationItem(confPrefix, "accessKeyId"))),
		WithAccessKeySecret(config.GetString(utils.GetConfigurationItem(confPrefix, "accessKeySecret"))),
		WithEndpoint(config.GetString(utils.GetConfigurationItem(confPrefix, "endpoint"))),
		WithSignName(config.GetString(utils.GetConfigurationItem(confPrefix, "signName"))),
		WithTempCode(config.GetString(utils.GetConfigurationItem(confPrefix, "tempCode"))),
	)
}

// Send 发送验证码参数 https://next.api.aliyun.com/api/Dysmsapi/2017-05-25/SendSms?lang=GO&params={}
func (ali *Adapter) Send(params *sms.Config) error {
	var (
		now  = time.Now().Unix()
		task = &sms.SendTask{
			CountryCode: params.CountryCode,
			Phone:       params.Phone,
			TempParams:  params.TempParams,
			Status:      "pending",
			Err:         "",
			SendAt:      now,
		}
	)
	if len(params.Code) == 0 {
		params.Code = utils.RandomSMSCode(6)
	}
	task.Code = params.Code
	if !ali.lock.Lock(sms.GetPhoneSendLockKey(params.Phone), time.Minute) {
		return fmt.Errorf("请勿频繁发送短信验证码")
	}
	if params.ExpireAt == 0 {
		params.ExpireAt = 300
	}
	task.ExpireAt = params.ExpireAt
	// 模拟发送短信验证码, 用于测试
	if ali.options.Analog {
		task.Status = "success"
		if err := ali.saveTask(params.Phone, task); err != nil {
			logger.Errorf("保存发送任务失败: %s", err.Error())
		}
		return nil
	}
	sendResp, err := ali.Client.SendSms(&client.SendSmsRequest{
		SignName:      &ali.options.SignName,
		TemplateCode:  &ali.options.TempCode,
		PhoneNumbers:  &params.Phone,
		TemplateParam: &params.TempParams,
	})
	if err != nil {
		return err
	}
	if *sendResp.Body.Code != "OK" {
		task.Status = "failed"
		task.Err = *sendResp.Body.Message
		if err := ali.saveTask(params.Phone, task); err != nil {
			logger.Errorf("保存发送任务失败: %s", err.Error())
		}
		_ = ali.lock.Unlock(sms.GetPhoneSendLockKey(params.Phone))
		return fmt.Errorf("发送短信失败: %s", *sendResp.Body.Message)
	}
	task.Status = "success"
	if err := ali.saveTask(params.Phone, task); err != nil {
		logger.Errorf("保存发送任务失败: %s", err.Error())
	}
	return nil
}

func (ali *Adapter) saveTask(phone string, task *sms.SendTask) error {
	taskJson, err := json.Marshal(task)
	if err != nil {
		return err
	}
	return ali.db.Set(
		sms.GetPhoneSendTaskKey(phone),
		string(taskJson),
		time.Duration(task.ExpireAt)*time.Second,
	).Err()
}

func (ali *Adapter) getTask(phone string) (*sms.SendTask, error) {
	task := &sms.SendTask{}
	taskJson, err := ali.db.Get(sms.GetPhoneSendTaskKey(phone)).Bytes()
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(taskJson, task); err != nil {
		return nil, err
	}
	return task, nil
}

// Verify 验证验证码
func (ali *Adapter) Verify(phone, code string) error {
	exist, err := ali.db.Exists(sms.GetPhoneSendTaskKey(phone)).Result()
	if err != nil || exist == 0 {
		return fmt.Errorf("验证码已失效")
	}
	task, err := ali.getTask(phone)
	if err != nil {
		return err
	}
	if task.Code != code {
		return fmt.Errorf("验证码错误")
	}
	if task.Status != "success" {
		return fmt.Errorf("验证码已失效")
	}
	return nil
}

func NewAliCloudAdapter(opts ...AliOption) sms.SMS {
	var options AliOptions
	for _, o := range opts {
		o(&options)
	}
	ali := &Adapter{
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
