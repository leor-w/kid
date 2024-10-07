package oauth2

type User struct {
	UserId   string `json:"sub"`            // 第三方平台的用户唯一标识
	Email    string `json:"email"`          // 第三方平台的用户邮箱
	EmailVer bool   `json:"email_verified"` // 第三方平台的用户邮箱是否已验证
	UserName string `json:"name"`           // 第三方平台的用户昵称
	Picture  string `json:"picture"`        // 第三方平台的用户头像
	Locale   string `json:"locale"`         // 第三方平台的用户地区
}
