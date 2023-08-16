package dingtalk

import "time"

type WebhookOptions struct {
	Robots []*WebhookRobotOptions
}

func WithWebhookRobot(robots []*WebhookRobotOptions) WebhookOption {
	return func(o *WebhookOptions) {
		o.Robots = robots
	}
}

type WebhookRobotOptions struct {
	Name    string
	Webhook string
	//Keyword    string
	SignSecret string
	//Security    string
	Security    bool
	ConnTimeout time.Duration
}

const (
	SecuritySign    = "sign"
	SecurityKeyword = "keyword"
)

func WithWebhookRobotName(name string) WebhookRobotOption {
	return func(o *WebhookRobotOptions) {
		o.Name = name
	}
}

func WithWebhook(webhook string) WebhookRobotOption {
	return func(o *WebhookRobotOptions) {
		o.Webhook = webhook
	}
}

//
//func WithKeyword(keyword string) WebhookRobotOption {
//	return func(o *WebhookRobotOptions) {
//		o.Keyword = keyword
//	}
//}

func WithSignSecret(signSecret string) WebhookRobotOption {
	return func(o *WebhookRobotOptions) {
		o.SignSecret = signSecret
	}
}

//func WithSecurity(security string) WebhookRobotOption {
//	return func(o *WebhookRobotOptions) {
//		o.Security = security
//	}
//}
