package dingtalk

type SendMessageRequest struct {
	RobotName string
	MsgKey    string
	MsgParam  string
}

type WithdrawMessageRequest struct {
	RobotName        string
	ProcessQueryKeys []string
}

type WebhookSendMessageReq struct {
	Name string
	Body map[string]interface{}
	//MsgType string
	//Content string
	//At      *WebhookSendMessageReqAt
}

type WebhookSendMessageReqAt struct {
	AtMobiles []string
	AtUserIds []string
	IsAtAll   bool
}
