package facebook

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"golang.org/x/oauth2"

	"github.com/lsjhtang/kid/config"
	"github.com/lsjhtang/kid/database/redis"
	"github.com/lsjhtang/kid/injector"
	"github.com/lsjhtang/kid/logger"
	plugin "github.com/lsjhtang/kid/plugin/oauth2"
	"github.com/lsjhtang/kid/utils"
)

type OAuth struct {
	options        Options
	facebookConfig oauth2.Config

	rds *redis.Client `inject:""`
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

// HandleOAuth2ByAuthCode 处理授权码登录授权
func (oauth *OAuth) HandleOAuth2ByAuthCode(code *plugin.VerifyCode) (*plugin.OAuthUser, error) {
	if exist := oauth.rds.Exists(GetOAuthURLIdentifierKey(code.State)); exist.Val() == 0 {
		return nil, fmt.Errorf("facebook oauth2: 授权来源未知或授权超时，请重新授权")
	}
	// 交换用户授权的 access_token
	token, err := oauth.facebookConfig.Exchange(context.Background(), code.Code)
	if err != nil {
		return nil, fmt.Errorf("facebook oauth2: 交换授权码失败: %s", err.Error())
	}

	return oauth.getUserInfoByAccessToken(token.AccessToken)
}

// HandleOAuth2ByAPPAuthToken 处理 APP 授权登录
func (oauth *OAuth) HandleOAuth2ByAPPAuthToken(token string) (*plugin.OAuthUser, error) {
	return oauth.getUserInfoByAccessToken(token)
}

// BuildAuthPageURL 构建授权页面 URL
func (oauth *OAuth) BuildAuthPageURL() (string, error) {
	state := utils.UUID()
	url := oauth.facebookConfig.AuthCodeURL(state)
	if err := oauth.rds.Set(GetOAuthURLIdentifierKey(state), url, time.Minute*30).Err(); err != nil {
		return "", fmt.Errorf("facebook oauth2: 设置 state 失败: %s", err.Error())
	}
	return url, nil
}

// getUserInfoByAccessToken 通过 access_token 获取用户信息
func (oauth *OAuth) getUserInfoByAccessToken(accessToken string) (*plugin.OAuthUser, error) {
	// 构建获取用户信息的请求地址
	userInfoUrl := fmt.Sprintf("%s?fields=%s&access_token=%s",
		EndpointUserInfo,
		strings.Join(oauth.options.Fields, ","),
		accessToken)
	// 请求 Facebook Graph API 获取用户信息
	resp, err := http.Get(userInfoUrl)
	if err != nil {
		return nil, fmt.Errorf("facebook oauth2: 获取用户信息失败: %s", err.Error())
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			logger.Errorf("facebook oauth2: 关闭获取信息响应流失败: %s", err.Error())
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("facebook oauth2: 获取用户信息失败: %s", string(body))
	}

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("facebook oauth2: 解析用户信息失败: %s", err.Error())
	}
	return &plugin.OAuthUser{
		UserId:   user.ID,
		Email:    user.Email,
		UserName: user.Name,
	}, nil
}

func New(opts ...Option) *OAuth {
	o := Options{}
	for _, opt := range opts {
		opt(&o)
	}
	o.Fields = append(o.Fields, "id", "email")
	return &OAuth{
		options: o,
		facebookConfig: oauth2.Config{
			ClientID:     o.ClientID,
			ClientSecret: o.ClientSecret,
			RedirectURL:  o.RedirectURL,
			Scopes:       o.Scope,
			Endpoint: oauth2.Endpoint{
				AuthURL:  EndpointAuthURL,
				TokenURL: EndpointTokenURL,
			},
		},
	}
}
