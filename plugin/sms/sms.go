package sms

import "fmt"

type Config struct {
	CountryCode string // 国家码, 如: +86, +1 等, 支持国际短信发送的平台需要
	Phone       string // 手机号
	IsAnalog    bool   // 本次是否模拟发送, 用于测试
	Code        string // 验证码, 需要发送的验证码, 可以为空, 由套件自动生成
	Language    string // 语言, 默认 en, 可以为空, 由套件自动生成
	Temp        string // 短信模板内容，这个是短信内容的模板，如: "您的验证码是: {code}" 等
	TempParams  string // 模板参数, json 格式, 如: {"code": "123456"} 或者 {"code": "123456", "product": "xxx"} 等。某些平台可能不需要
	ExpireAt    int    // 过期时间, 默认 300 秒（5分钟）。可以通过此参数自定义设置验证码的有效时间
}

type SendTask struct {
	CountryCode string `json:"country_code"` // 国家码
	Phone       string `json:"phone"`        // 手机号
	Code        string `json:"code"`         // 验证码
	TempParams  string `json:"temp_params"`  // 模板参数 json 格式
	Status      string `json:"status"`       // 发送状态
	Err         string `json:"err"`          // 错误信息
	SendAt      int64  `json:"send_at"`      // 发送时间
	ExpireAt    int    `json:"expire_at"`    // 过期时间
}

type SMS interface {
	Send(params *Config) error       // 发送短信验证码
	Verify(phone, code string) error // 验证短信验证码
}

const (
	PhoneSendTaskKey = "phone:sms:send:task.%s" // 手机号发送任务 redis key
	PhoneSendLockKey = "phone:sms:send:lock.%s" // 手机号发送锁 redis key (防止重复发送)
)

func GetPhoneSendTaskKey(phone string) string {
	return fmt.Sprintf(PhoneSendTaskKey, phone)
}

func GetPhoneSendLockKey(phone string) string {
	return fmt.Sprintf(PhoneSendLockKey, phone)
}
