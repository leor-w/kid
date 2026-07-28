package instagram

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/instagram"

	"github.com/lsjhtang/kid/config"
	"github.com/lsjhtang/kid/database/redis"
	"github.com/lsjhtang/kid/injector"
	"github.com/lsjhtang/kid/logger"
	"github.com/lsjhtang/kid/plugin/lock"
	plugin "github.com/lsjhtang/kid/plugin/oauth2"
	"github.com/lsjhtang/kid/utils"
)

type OAuth struct {
	options       Options
	instagramConf oauth2.Config
	rds           *redis.Client `inject:""`
	lock          lock.Lock     `inject:""`
}

func (oauth *OAuth) Provide(ctx context.Context) any {
	var confName string
	name, ok := ctx.Value(injector.NameKey{}).(string)
	if ok && len(name) > 0 {
		confName = "." + name
	}
	confPrefix := fmt.Sprintf("oauth2%s", confName)
	if !config.Exist(confPrefix) {
		panic(fmt.Sprintf("配置文件为找到 [%s.*]，请检查配置文件", confPrefix))
	}
	return New(
		WithClientID(config.GetString(utils.GetConfigurationItem(confPrefix, "client_id"))),
		WithClientSecret(config.GetString(utils.GetConfigurationItem(confPrefix, "client_secret"))),
		WithRedirectURL(config.GetString(utils.GetConfigurationItem(confPrefix, "redirect_url"))),
		WithScope(config.GetStringSlice(utils.GetConfigurationItem(confPrefix, "scope"))),
	)
}

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Option func(o *Options)

func (oauth *OAuth) HandleOAuth2ByAuthCode(code *plugin.VerifyCode) (*plugin.OAuthUser, error) {
	if !oauth.stateExist(code.State) {
		return nil, fmt.Errorf("instagram oauth2: 授权来源未知或授权超时，请重新授权")
	}
	token, err := oauth.ExchangeAccessToken(code.Code)
	if err != nil {
		return nil, fmt.Errorf("instagram oauth2: 交换授权码失败: %s", err.Error())
	}
	// 请求 Facebook Graph API 获取用户信息
	resp, err := oauth.GetUserInfoByAccessToken(token)
	if err != nil {
		return nil, err
	}

	return &plugin.OAuthUser{
		UserId:   resp.ID,
		UserName: resp.UserName,
	}, nil
}

// ExchangeAccessToken 交换用户授权的 access_token
func (oauth *OAuth) ExchangeAccessToken(code string) (string, error) {
	resp, err := http.PostForm(EndpointExchangeAccessTokenURL, url.Values{
		"client_id":     {oauth.options.ClientID},
		"client_secret": {oauth.options.ClientSecret},
		"grant_type":    {"authorization_code"},
		"redirect_uri":  {oauth.options.RedirectURL},
		"code":          {code},
	})
	if err != nil {
		return "", fmt.Errorf("instagram oauth2: 交换 access_token 失败: %s", err.Error())
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			logger.Errorf("instagram oauth2: 关闭请求响应失败: %s", err.Error())
		}
	}(resp.Body)
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("instagram oauth2: 交换 access_token 失败: %s", string(body))
	}
	var exchangeResp ExchangeAccessTokenResp
	if err := json.NewDecoder(resp.Body).Decode(&exchangeResp); err != nil {
		return "", fmt.Errorf("instagram oauth2: 解析 access_token 失败: %s", err.Error())
	}
	return exchangeResp.AccessToken, nil
}

// GetUserInfoByAccessToken 通过 access_token 获取用户信息
func (oauth *OAuth) GetUserInfoByAccessToken(accessToken string) (*GetUserResp, error) {
	userInfoURL := fmt.Sprintf("%s?fields=%s&access_token=%s",
		EndpointUserInfoURL,
		strings.Join(oauth.options.Scope, ","),
		accessToken)
	resp, err := http.Get(userInfoURL)
	if err != nil {
		return nil, fmt.Errorf("instagram oauth2: 获取用户信息失败: %s", err.Error())
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			logger.Errorf("instagram oauth2: 关闭获取用户信息响应流失败: %s", err.Error())
		}
	}(resp.Body)
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("instagram oauth2: 获取用户信息失败: %s", string(body))
	}

	var user GetUserResp
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("instagram oauth2: 解析用户信息失败: %s", err.Error())
	}
	return &user, nil
}

// HandleOAuth2ByAPPAuthToken 处理 APP 授权登录
func (oauth *OAuth) HandleOAuth2ByAPPAuthToken(token string) (*plugin.OAuthUser, error) {
	resp, err := oauth.GetUserInfoByAccessToken(token)
	if err != nil {
		return nil, err
	}
	return &plugin.OAuthUser{
		UserId:   resp.ID,
		UserName: resp.UserName,
	}, nil
}

// stateExist 判断 state 是否存在
func (oauth *OAuth) stateExist(state string) bool {
	return oauth.rds.Exists(GetOAuthURLIdentifierKey(state)).Val() > 0
}

// BuildAuthPageURL 构建授权页面 URL
func (oauth *OAuth) BuildAuthPageURL() (string, error) {
	state := utils.UUID()
	authUrl := oauth.instagramConf.AuthCodeURL(state)
	if err := oauth.rds.Set(GetOAuthURLIdentifierKey(state), authUrl, 30*time.Minute).Err(); err != nil {
		return "", fmt.Errorf("instagram oauth2: 设置 state 失败: %s", err.Error())
	}
	return authUrl, nil
}

func New(opts ...Option) *OAuth {
	o := Options{}
	for _, opt := range opts {
		opt(&o)
	}
	return &OAuth{
		options: o,
		instagramConf: oauth2.Config{
			ClientID:     o.ClientID,
			ClientSecret: o.ClientSecret,
			RedirectURL:  o.RedirectURL,
			Scopes:       o.Scope,
			Endpoint:     instagram.Endpoint,
		},
	}
}
