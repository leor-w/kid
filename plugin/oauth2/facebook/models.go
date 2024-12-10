package facebook

import "fmt"

const (
	OAuthUrlIdentifierKey = "facebook.oauth2.url.identifier.%s"
)

func GetOAuthURLIdentifierKey(state string) string {
	return fmt.Sprintf(OAuthUrlIdentifierKey, state)
}

const (
	// EndpointAuthURL 构建用户授权地址的请求地址
	EndpointAuthURL = "https://www.facebook.com/v21.0/dialog/oauth"
	// EndpointTokenURL 交换用户授权的 access_token 的请求地址
	EndpointTokenURL = "https://graph.facebook.com/v21.0/oauth/access_token"
	// EndpointUserInfo 获取用户信息的请求地址
	EndpointUserInfo = "https://graph.facebook.com/me"
)

const (
	ScopeEmail = "email"
)
