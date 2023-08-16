package dingtalk

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/imroc/req/v3"
)

type DingtalkTest struct {
	Dingtalk *Webhook `inject:""`
}

func (d *DingtalkTest) Provide(_ context.Context) interface{} {
	return d
}

func TestDingtalk(t *testing.T) {
	var message = map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]interface{}{
			"title": "notify：",
			"text":  "我会倒立嘘嘘",
		},
	}
	var msg, _ = json.Marshal(message)
	var resp, err = req.
		SetHeader("Content-Type", "application/json").
		SetBody(msg).
		Post("https://oapi.dingtalk.com/robot/send?access_token=2275ca9fb9933b1c37c653f4e583ba0703ea479752bb16e9ec72f3019a85c1f7")
	bodyData, err := resp.ToBytes()
	fmt.Println(string(bodyData))
	if err != nil {
		t.Error("ExceptionHandler.Report: 机器人通知失败")
		t.Fail()
	}
}

func TestWebhook(t *testing.T) {
	dingtalk := NewWebhookRobot(&WebhookRobotOptions{
		Name:        "测试机器人",
		Webhook:     "在这里输入机器人的 webhook",
		SignSecret:  "在这里填入加签秘钥",
		Security:    true, // 是否使用加签
		ConnTimeout: 3,
	})

	if err := dingtalk.PostMessage(map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]interface{}{
			"title": "测试消息",
			"text":  "### 你可以教我么",
		},
	}); err != nil {
		t.Error(err.Error())
		t.Fail()
	}
}
