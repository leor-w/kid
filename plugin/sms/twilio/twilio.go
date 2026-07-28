package twilio

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/lsjhtang/kid/database/redis"
	"github.com/lsjhtang/kid/plugin/lock"

	"github.com/lsjhtang/kid/logger"

	"github.com/lsjhtang/kid/plugin/sms"

	"github.com/lsjhtang/kid/config"
	"github.com/lsjhtang/kid/plugin"
	"github.com/lsjhtang/kid/utils"

	"github.com/twilio/twilio-go"
	verify "github.com/twilio/twilio-go/rest/verify/v2"
)

type Option func(*Options)

type Adapter struct {
	client  *twilio.RestClient
	options *Options
	lock    lock.Lock     `inject:""`
	db      *redis.Client `inject:""`
}

func (t *Adapter) Provide(ctx context.Context) any {
	var configName string
	name, ok := ctx.Value(new(plugin.NameKey)).(string)
	if ok && len(name) > 0 {
		configName = "." + name
	}
	confPrefix := fmt.Sprintf("twilio%s", configName)
	if !config.Exist(fmt.Sprintf(confPrefix)) {
		panic(fmt.Sprintf("配置文件未找到配置项 [%s]", confPrefix))
	}
	return NewTwilio(
		WithAnalog(config.GetBool(utils.GetConfigurationItem(confPrefix, "analog"))),
		WithAccountSid(config.GetString(utils.GetConfigurationItem(confPrefix, "account_sid"))),
		WithAuthToken(config.GetString(utils.GetConfigurationItem(confPrefix, "auth_token"))),
		WithServiceSid(config.GetString(utils.GetConfigurationItem(confPrefix, "service_sid"))),
		WithEdge(config.GetString(utils.GetConfigurationItem(confPrefix, "edge"))),
		WithRegion(config.GetString(utils.GetConfigurationItem(confPrefix, "region"))),
	)
}

func (t *Adapter) Send(params *sms.Config) error {
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
	if !t.lock.Lock(sms.GetPhoneSendLockKey(params.Phone), time.Minute) {
		return fmt.Errorf("请勿频繁发送短信验证码")
	}
	if params.ExpireAt == 0 {
		params.ExpireAt = 300
	}
	task.ExpireAt = params.ExpireAt
	// 模拟发送短信验证码, 用于测试
	if t.options.Analog || params.IsAnalog {
		task.Code = "888888"
		task.Status = "success"
		if err := t.saveTask(params.Phone, task); err != nil {
			logger.Errorf("保存发送任务失败: %s", err.Error())
		}
		_ = t.lock.Unlock(sms.GetPhoneSendLockKey(params.Phone))
		return nil
	}
	smsCode := &verify.CreateVerificationParams{}
	smsCode.SetTo(params.Phone)
	smsCode.SetChannel("sms")
	if params.Temp != "" {
		if strings.Contains(params.Temp, "{code}") {
			params.Temp = strings.Replace(params.Temp, "{code}", params.Code, -1)
		} else {
			params.Temp += " " + params.Code
		}
		smsCode.SetCustomMessage(params.Temp)
	} else {
		smsCode.SetCustomCode(params.Code)
		smsCode.SetLocale(params.Language)
		if params.TempParams != "" {
			smsCode.SetTemplateCustomSubstitutions(params.TempParams)
		}
	}
	resp, err := t.client.VerifyV2.CreateVerification(t.options.ServiceSid, smsCode)
	if err != nil {
		return fmt.Errorf("发送短信验证码失败: %s", err.Error())
	}
	if *resp.Status != "pending" {
		return fmt.Errorf("发送短信验证码失败: %s", *resp.Status)
	}
	if *resp.Status != "pending" {
		task.Status = "failed"
		task.Err = *resp.Status
		if err := t.saveTask(params.Phone, task); err != nil {
			logger.Errorf("保存发送任务失败: %s", err.Error())
		}
		_ = t.lock.Unlock(sms.GetPhoneSendLockKey(params.Phone))
		return fmt.Errorf("发送短信失败: %s", *resp.Status)
	}
	task.Status = "success"
	if err := t.saveTask(params.Phone, task); err != nil {
		logger.Errorf("保存发送任务失败: %s", err.Error())
	}

	return err
}

func (t *Adapter) saveTask(phone string, task *sms.SendTask) error {
	taskJson, err := json.Marshal(task)
	if err != nil {
		return err
	}
	return t.db.Set(
		sms.GetPhoneSendTaskKey(phone),
		string(taskJson),
		time.Duration(task.ExpireAt)*time.Second,
	).Err()
}

func (t *Adapter) getTask(phone string) (*sms.SendTask, error) {
	task := &sms.SendTask{}
	taskJson, err := t.db.Get(sms.GetPhoneSendTaskKey(phone)).Bytes()
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(taskJson, task); err != nil {
		return nil, err
	}
	return task, nil
}

func (t *Adapter) Verify(phone, code string) error {
	exist, err := t.db.Exists(sms.GetPhoneSendTaskKey(phone)).Result()
	if err != nil || exist == 0 {
		return fmt.Errorf("验证码已失效")
	}
	task, err := t.getTask(phone)
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

func NewTwilio(options ...Option) *Adapter {
	opt := &Options{}
	for _, o := range options {
		o(opt)
	}
	err := os.Setenv("TWILIO_ACCOUNT_SID", opt.AccountSid)
	err = os.Setenv("TWILIO_AUTH_TOKEN", opt.AuthToken)
	if err != nil {
		panic(fmt.Sprintf("Adapter 初始化错误: %s", err.Error()))
	}
	return &Adapter{
		client:  twilio.NewRestClient(),
		options: opt,
	}
}
