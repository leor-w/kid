package qq

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/lsjhtang/kid/config"
	"github.com/lsjhtang/kid/plugin"

	"github.com/lsjhtang/kid/utils"
)

const (
	URL = "https://graph.qq.com/oauth2.0/"
)

type OAuth struct {
	options *Options
}

func (auth *OAuth) Provide(ctx context.Context) interface{} {
	var configName string
	name, ok := ctx.Value(new(plugin.NameKey)).(string)
	if ok && len(name) > 0 {
		configName = "." + name
	}
	confPrefix := fmt.Sprintf("qq%s", configName)
	if !config.Exist(fmt.Sprintf(confPrefix)) {
		panic(fmt.Sprintf("config.yaml file not found configuration item [%s]", confPrefix))
	}
	return New(
		WithAppId(config.GetString(utils.GetConfigurationItem(confPrefix, "oauth.appId"))),
		WithAppKey(config.GetString(utils.GetConfigurationItem(confPrefix, "oauth.appKey"))),
		WithDisplay(config.GetString(utils.GetConfigurationItem(confPrefix, "oauth.display"))),
		WithRedirectUri(config.GetString(utils.GetConfigurationItem(confPrefix, "oauth.redirectUri"))),
		WithProduct(config.GetBool(utils.GetConfigurationItem(confPrefix, "oauth.product"))),
	)
}

// GetAuthCode 获取登录的 Auth code 返回登录扫码的 URL
func (auth *OAuth) GetAuthCode(state string) (string, error) {
	var params utils.ReqParams
	params["response_type"] = "code"
	params["client_id"] = auth.options.AppId
	params["redirect_uri"] = auth.options.RedirectUri
	params["state"] = state
	params["scope"] = ""
	params["display"] = auth.options.Display
	resp, err := utils.GET(URL+"authorize", nil, params)
	if err != nil {
		return "", err
	}
	var authCode AuthCodeResp
	if err := json.Unmarshal(resp, &authCode); err != nil {
		return "", err
	}
	return authCode.Code, nil
}

// GetToken 通过 code 获取 accessToken
func (auth *OAuth) GetToken(code string) (*AccessTokenResp, error) {
	params := utils.ReqParams{}
	params["grant_type"] = "authorization_code"
	params["client_id"] = auth.options.AppId
	params["client_secret"] = auth.options.AppKey
	params["redirect_uri"] = auth.options.RedirectUri
	params["code"] = code
	resp, err := utils.GET(URL+"token", nil, params)
	if err != nil {
		return nil, err
	}
	var info AccessTokenResp
	if err := json.Unmarshal(resp, &info); err != nil {
		return nil, errors.New(fmt.Sprintf("OAuth.GetToken failed: %s", err.Error()))
	}
	return &info, nil
}

// GetOpenId 通过 accessToken 换取用户 openId
func (auth *OAuth) GetOpenId(accessToken string) (string, error) {
	params := utils.ReqParams{}
	params["access_token"] = accessToken
	params["fmt"] = "json"
	resp, err := utils.GET(URL+"me", nil, params)
	if err != nil {
		return "", err
	}
	var openid OpenidResp
	if err := json.Unmarshal(resp, &openid); err != nil {
		return "", err
	}
	return openid.Openid, nil
}

// GetOpenIdWithAuthCode 获取 auth code 直接获取用户 openid
func (auth *OAuth) GetOpenIdWithAuthCode(authCode string) (string, error) {
	info, err := auth.GetToken(authCode)
	if err != nil {
		return "", err
	}
	openId, err := auth.GetOpenId(info.AccessToken)
	if err != nil {
		return "", err
	}
	return openId, nil
}

func (auth *OAuth) GetAppId() string {
	return auth.options.AppId
}

func New(opts ...Option) *OAuth {
	var options Options
	for _, opt := range opts {
		opt(&options)
	}
	return &OAuth{options: &options}
}
