package qq

type AccessTokenResp struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    string `json:"expires_in"`
}

type AuthCodeResp struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}

type OpenidResp struct {
	ClientId string `json:"client_id"`
	Openid   string `json:"openid"`
}
