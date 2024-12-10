package instagram

import "fmt"

const (
	// EndpointExchangeAccessTokenURL 交换用户授权的 access_token 的请求地址
	EndpointExchangeAccessTokenURL = "https://api.instagram.com/oauth/access_token"
	// EndpointUserInfoURL 获取用户信息的请求地址
	EndpointUserInfoURL = "https://graph.instagram.com/me"
)

const (
	OAuthUrlIdentifierKey = "instagram.oauth2.url.identifier.%s"
)

// GetOAuthURLIdentifierKey 获取 OAuth URL 标识键
func GetOAuthURLIdentifierKey(state string) string {
	return fmt.Sprintf(OAuthUrlIdentifierKey, state)
}

type ExchangeAccessTokenResp struct {
	AccessToken string `json:"access_token"`
	UserID      string `json:"user_id"`
}

type GetUserResp struct {
	ID       string `json:"id"`
	UserName string `json:"username"`
}
