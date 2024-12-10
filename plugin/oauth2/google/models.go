package google

import "fmt"

const (
	OAuthUrlIdentifierKey = "google.oauth2.url.identifier.%s"
)

func GetOAuthURLIdentifierKey(state string) string {
	return fmt.Sprintf(OAuthUrlIdentifierKey, state)
}

const (
	// EndpointUserInfo 获取用户信息的请求地址
	EndpointUserInfo = "https://www.googleapis.com/oauth2/v3/userinfo"
)

const (
	// ScopeProfile Google 获取用户信息的权限
	ScopeProfile = "https://www.googleapis.com/auth/userinfo.profile" // 获取用户信息
	// ScopeEmail Google 获取用户邮箱的权限
	ScopeEmail = "https://www.googleapis.com/auth/userinfo.email" // 获取用户邮箱
)
