package oauth2

type User struct {
	UserId   string `json:"sub"`            // 第三方平台的用户唯一标识
	Email    string `json:"email"`          // 第三方平台的用户邮箱
	EmailVer bool   `json:"email_verified"` // 第三方平台的用户邮箱是否已验证
	UserName string `json:"name"`           // 第三方平台的用户昵称
	Picture  string `json:"picture"`        // 第三方平台的用户头像
	Locale   string `json:"locale"`         // 第三方平台的用户地区
}

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
}

type OAuth2 interface {
	// GetAuthPageURL 获取授权页面 URL
	GetAuthPageURL() string
	// HandleAuth 处理授权码
	HandleAuth(code *VerifyCode) (*User, error)
}
