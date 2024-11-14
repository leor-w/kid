package tiktok

const (
	endpointAuth     = "https://www.tiktok.com/v2/auth/authorize/"    // 授权页面地址
	endpointToken    = "https://open.tiktokapis.com/v2/oauth/token/"  // 获取、刷新 access_token 地址
	endpointRevoke   = "https://open.tiktokapis.com/v2/oauth/revoke/" // 注销 token 地址
	endpointUserInfo = "https://open.tiktokapis.com/v2/user/info/"    // 获取用户信息地址
)

const (
	UserScopeBasic        = "user.info.basic"
	UserScopeArtistRead   = "artist.certification.read"
	UserScopeArtistUpdate = "artist.certification.update"
	UserScopeUserProfile  = "user.info.profile"
	UserScopeUserStats    = "user.info.stats"
)

// FetchAccessTokenReq 获取 access_token 请求参数
type FetchAccessTokenReq struct {
	/* required 客户端 Key 在tiktok开发者平台查看 */
	ClientKey string `json:"client_key"`
	/* required 客户端 Secret 在tiktok开发者平台查看 */
	ClientSecret string `json:"client_secret"`
	/* required 授权码 */
	Code string `json:"code"`
	/* required 授权类型 固定为：authorization_code */
	GrantType string `json:"grant_type"`
	/* required 回调的重定向地址 */
	RedirectURL string `json:"redirect_uri"`
	/* 仅适用于移动和桌面应用程序。代码验证器用于在 PKCE 授权流程中生成代码质询。 */
	CodeVerifier string `json:"code_verifier"`
}

// FetchAccessTokenResp 获取 access_token 响应参数
type FetchAccessTokenResp struct {
	/* 用户授权的唯一标识 */
	OpenId string `json:"open_id"`
	/* 用户授权的数据范围 */
	Scope string `json:"scope"`
	/* 用户授权的访问令牌 */
	AccessToken string `json:"access_token"`
	/* 访问令牌的有效期，单位为秒 */
	ExpiresIn int64 `json:"expires_in"`
	/* 用户授权的刷新令牌, 首次前发后365天有效 */
	RefreshToken string `json:"refresh_token"`
	/* 刷新令牌的有效期，单位为秒 */
	RefreshExpiresIn int64 `json:"refresh_expires_in"`
	/* 用户授权的令牌类型 固定为：Bearer */
	TokenType string `json:"token_type"`

	/* 以下为请求错误响应内容 */
	/* 错误类型 */
	Error string `json:"error"`
	/* 错误描述 */
	ErrorDescription string `json:"error_description"`
	/* LogId */
	LogId string `json:"log_id"`
}

type UserInfoResp struct {
	/* 用户信息 */
	Data map[string]UserObject `json:"data"`
	/* 错误信息 */
	Error Error `json:"error"`
}

type Error struct {
	/* 错误类别 */
	Code string `json:"code"`
	/* 错误描述 */
	Message string `json:"message"`
	/* LogId */
	LogId string `json:"log_id"`
}

type UserObject struct {
	/* 用户在当前应用中的唯一标识 */
	OpenId string `json:"open_id"`
	/* 用户在当前开发者账号下所有应用的唯一标识 */
	UnionId string `json:"union_id"`
	/* 用户头像地址 */
	AvatarURL string `json:"avatar_url"`
	/* 用户头像地址（100x100） */
	AvatarURL100 string `json:"avatar_url_100"`
	/* 最高分辨率的用户头像地址 */
	AvatarLargeURL string `json:"avatar_large_url"`
	/* 显示名称 */
	DisplayName string `json:"display_name"`
	/* 用户个人简介 */
	BioDescription string `json:"bio_description"`
	/* 用户TikTok个人资料页面链接 */
	ProfileDeepLink string `json:"profile_deep_link"`
	/* 用户是否已完成TikTok验证 */
	IsVerified bool `json:"is_verified"`
	/* 用户昵称 */
	Username string `json:"username"`
	/* 用户的关注者数量 */
	FollowerCount int64 `json:"follower_count"`
	/* 用户的关注的账户数量 */
	FollowingCount int64 `json:"following_count"`
	/* 用户所有视频的点赞数量 */
	LikesCount int64 `json:"likes_count"`
	/* 用户公开发布的视频总数 */
	VideoCount int64 `json:"video_count"`
}
