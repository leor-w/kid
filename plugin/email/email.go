package email

import (
	"context"
	"fmt"

	"github.com/lsjhtang/kid/injector"

	"github.com/lsjhtang/kid/config"
	"github.com/lsjhtang/kid/utils"

	"gopkg.in/gomail.v2"
)

type Email struct {
	dialer  *gomail.Dialer
	options *Options
}

func (e *Email) Provide(ctx context.Context) interface{} {
	var confName string
	if name, ok := ctx.Value(injector.NameKey{}).(string); ok && len(name) > 0 {
		confName = "." + name
	}
	confPrefix := fmt.Sprintf("email%s", confName)
	if !config.Exist(confPrefix) {
		panic(fmt.Sprintf("config.yaml file not found configuration item [%s]", confPrefix))
	}
	return New(
		WithHost(config.GetString(utils.GetConfigurationItem(confPrefix, "host"))),
		WithPort(config.GetInt(utils.GetConfigurationItem(confPrefix, "port"))),
		WithPwd(config.GetString(utils.GetConfigurationItem(confPrefix, "pwd"))),
		WithFrom(config.GetString(utils.GetConfigurationItem(confPrefix, "from"))),
	)
}

func (e *Email) Send(msg *Message) error {
	var waitSend []*gomail.Message
	for _, toAddr := range msg.ToAddrs {
		var sendMsg = gomail.NewMessage()
		sendMsg.SetHeader("From", e.options.From)
		sendMsg.SetHeader("To", toAddr)
		for _, cc := range msg.Cc {
			sendMsg.SetHeader("Cc", cc)
		}
		sendMsg.SetHeader("Subject", msg.Subject)
		sendMsg.SetBody(msg.Body.MsgType, msg.Body.Body)
		if len(msg.Attach) > 0 {
			sendMsg.Attach(msg.Attach)
		}
		waitSend = append(waitSend, sendMsg)
	}
	return e.dialer.DialAndSend(waitSend...)
}

func (e *Email) Init() {
	dialer := gomail.NewDialer(e.options.Host, e.options.Port, e.options.From, e.options.Pwd)
	e.dialer = dialer
}

func New(opts ...Option) *Email {
	var options Options
	for _, opt := range opts {
		opt(&options)
	}
	var email = &Email{options: &options}
	email.Init()
	return email
}
