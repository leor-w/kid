package oauth2

type OAuthUser struct {
	UserId   string `json:"sub"`            // 第三方平台的用户唯一标识
	Email    string `json:"email"`          // 第三方平台的用户邮箱
	EmailVer bool   `json:"email_verified"` // 第三方平台的用户邮箱是否已验证
	UserName string `json:"name"`           // 第三方平台的用户昵称
	Picture  string `json:"picture"`        // 第三方平台的用户头像
	Locale   string `json:"locale"`         // 第三方平台的用户地区
}

const (
	PlatformTiktok    = "tiktok"    // tiktok
	PlatformApple     = "apple"     // apple
	PlatformFacebook  = "facebook"  // facebook
	PlatformGoogle    = "google"    // google
	PlatformInstagram = "instagram" // instagram
)

type VerifyCode struct {
	/*
		Code
		第三方平台返回的授权码，用于获取用户信息
	*/
	Code string
	/*
		CodeVerifier
		TikTok 在移动设备或桌面设备验证时必传参数，用于验证 code 的有效性
		具体参考：https://developers.tiktok.com/doc/oauth-user-access-token-management?enter_method=left_navigation
	*/
	CodeVerifier string
	/*
		Fields
		TikTok 获取用户信息时需要的字段，例如：open_id, username, avatar_url;
		TikTok 特有参数。
		可以不传，套件内已经默认集成了open_id, union_id, avatar_url这三个字段信息
		具体可选内容参考：https://developers.tiktok.com/doc/tiktok-api-v2-get-user-info?enter_method=left_navigation
	*/
	Fields []string

	/*
		State
		第三方平台依据 BuildAuthPageURL 设置的 state 参数进行原样返回的 state 参数，用于验证请求的合法性
	*/
	State string

	/*
		IdToken
		Apple 授权登录时返回的 id_token
	*/
	IdToken string
}

type OAuth2 interface {
	// BuildAuthPageURL 获取授权页面 URL
	BuildAuthPageURL() (string, error)

	// HandleOAuth2ByAuthCode 处理授权码登录授权
	HandleOAuth2ByAuthCode(*VerifyCode) (*OAuthUser, error)

	// HandleOAuth2ByAPPAuthToken 处理 APP 授权登录
	HandleOAuth2ByAPPAuthToken(string) (*OAuthUser, error)
}
