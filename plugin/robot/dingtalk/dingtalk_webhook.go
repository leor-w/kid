package dingtalk

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/tidwall/gjson"

	robot2 "github.com/lsjhtang/kid/plugin/robot"

	"github.com/lsjhtang/kid/config"
	"github.com/lsjhtang/kid/plugin"
	"github.com/lsjhtang/kid/utils"
	"github.com/spf13/cast"

	"github.com/imroc/req/v3"
)

type Webhook struct {
	robots map[string]*WebhookRobot
}

type WebhookOption func(*WebhookOptions)

func (hook *Webhook) Provide(ctx context.Context) interface{} {
	var configName string
	name, ok := ctx.Value(new(plugin.NameKey)).(string)
	if ok && len(name) > 0 {
		configName = "." + name
	}
	confPrefix := fmt.Sprintf("dingtalk.webhook%s", configName)
	if !config.Exist(fmt.Sprintf(confPrefix)) {
		panic(fmt.Sprintf("配置文件中未找到对应配置项 [%s], 请检查", confPrefix))
	}
	robotConfs, ok := config.Get(utils.GetConfigurationItem(confPrefix, "robots")).([]interface{})
	if !ok {
		panic(fmt.Sprintf("钉钉 Webhook 机器人配置错误, 请检查"))
	}
	var robots []*WebhookRobotOptions
	for _, conf := range robotConfs {
		robotConf, ok := conf.(map[string]interface{})
		if !ok {
			panic(fmt.Sprintf("钉钉 Webhook 机器人配置错误, 请检查"))
		}
		robots = append(robots, &WebhookRobotOptions{
			Name:        cast.ToString(robotConf["name"]),
			Webhook:     cast.ToString(robotConf["webhook"]),
			SignSecret:  cast.ToString(robotConf["sign_secret"]),
			Security:    cast.ToBool(robotConf["security"]),
			ConnTimeout: cast.ToDuration(robotConf["conn_timeout"]) * time.Second,
		})
	}
	return NewWebhook(WithWebhookRobot(robots))
}

func (hook *Webhook) SendMessage(params interface{}) error {
	reqConf, ok := params.(*WebhookSendMessageReq)
	if !ok {
		return robot2.ErrWrongParamType
	}
	sendRobot, exist := hook.robots[reqConf.Name]
	if !exist {
		return ErrNotFoundRobot
	}
	if err := sendRobot.SendMessage(reqConf); err != nil {
		return err
	}
	return nil
}

func (hook *Webhook) WithdrawMessage(params interface{}) error {
	return nil
}

func NewWebhook(opts ...WebhookOption) *Webhook {
	options := &WebhookOptions{}
	for _, opt := range opts {
		opt(options)
	}
	var robots = make(map[string]*WebhookRobot)
	for _, robotConf := range options.Robots {
		robots[robotConf.Name] = NewWebhookRobot(robotConf)
	}
	return &Webhook{
		robots: robots,
	}
}

type WebhookRobot struct {
	client  *req.Client
	options *WebhookRobotOptions
}

type WebhookRobotOption func(*WebhookRobotOptions)

func (robot *WebhookRobot) SendMessage(msg *WebhookSendMessageReq) error {
	return robot.PostMessage(msg.Body)
}

func (robot *WebhookRobot) PostMessage(body map[string]interface{}) error {
	msg, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("请求数据错误: %s", err.Error())
	}
	webhook := robot.options.Webhook
	if robot.options.Security {
		now := time.Now().UnixMilli()
		sign := robot.Sign(now)
		webhook = fmt.Sprintf("%s&timestamp=%d&sign=%s", webhook, now, sign)
	}
	resp, err := robot.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(msg).
		Post(webhook)
	if err != nil {
		return fmt.Errorf("请求错误: %s", err.Error())
	}
	if resp.IsError() {
		return fmt.Errorf("请求响应错误: %s", resp.Err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("请求错误, 错误码: %d", resp.StatusCode)
	}
	respBodyBytes, err := resp.ToBytes()
	if err != nil {
		return fmt.Errorf("读取响应内容错误: %s", err.Error())
	}
	respBody := gjson.ParseBytes(respBodyBytes)
	if respBody.Get("errcode").Int() != 0 {
		return fmt.Errorf("发送消息错误: 错误码 [%d], 错误描述 [%s]",
			respBody.Get("errcode").Int(),
			respBody.Get("errmsg").String())
	}
	return nil
}

func (robot *WebhookRobot) Sign(timestamp int64) string {
	hash := hmac.New(sha256.New, []byte(robot.options.SignSecret))
	hash.Write([]byte(fmt.Sprintf("%d\n%s", timestamp, robot.options.SignSecret)))
	return url.QueryEscape(base64.StdEncoding.EncodeToString(hash.Sum(nil)))
}

func NewWebhookRobot(options *WebhookRobotOptions) *WebhookRobot {
	cli := req.NewClient()
	if options.ConnTimeout < time.Second*3 {
		options.ConnTimeout = time.Second * 3
	}
	cli.SetTimeout(options.ConnTimeout)
	return &WebhookRobot{
		client:  cli,
		options: options,
	}
}
