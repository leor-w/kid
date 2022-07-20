package smscode

type SMSCode interface {
	Init() error
	SendSMS(phone, text string) error
}
