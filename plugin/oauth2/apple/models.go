package apple

import "fmt"

const (
	// BuildOAuthPageURL Apple OAuth2 构建授权页面 URL 的请求地址
	BuildOAuthPageURL = "https://appleid.apple.com/auth/authorize"
	// EndpointPublicKeyURL Apple OAuth2 获取公钥的请求地址
	EndpointPublicKeyURL = "https://appleid.apple.com/auth/keys"
)

const (
	// OAuthUrlIdentifierKey 构建生成的授权 URL 标识键，用于验证授权来源
	OAuthUrlIdentifierKey = "apple.oauth2.url.identifier.%s"
)

func GetOAuthURLIdentifierKey(state string) string {
	return fmt.Sprintf(OAuthUrlIdentifierKey, state)
}

// GetPublicKeyResp 获取公钥响应参数
type GetPublicKeyResp struct {
	Keys []struct {
		Kid string `json:"kid"`
		N   string `json:"n"`
		E   string `json:"e"`
	} `json:"keys"`
}
