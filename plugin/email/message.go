package email

type Message struct {
	ToAddrs []string
	Cc      []string
	Subject string
	Body    *Body
	Attach  string
}

type Body struct {
	MsgType string
	Body    string
}
