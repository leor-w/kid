package smscode

type SMSCode interface {
	SendSMS(phone, code string) error                                                                               // 发送短信验证码
	SendSMSWithParams(phone string, values map[string]interface{}) error                                            // 发送短信指定参数
	SendSMSWithTemplate(phone string, template string, params map[string]interface{}) error                         // 发送短信指定模版及参数
	SendSMSWithSignAndTemplate(phone string, signName string, template string, params map[string]interface{}) error // 发送短信指定短信签名、模版及参数
}
