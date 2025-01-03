package wechat

import "fmt"

const (
	ScopeBase     = "snsapi_base"     // 静默授权
	ScopeUserInfo = "snsapi_userinfo" // 用户信息授权
)

const (
	wechatOauthState = "wechat.oauth.state.%s"
)

func getWechatOauthStateKey(state string) string {
	return fmt.Sprintf(wechatOauthState, state)
}
