package tiktok

import (
	"context"
	"encoding/json"
	"fmt"

	localOauth2 "github.com/leor-w/kid/plugin/oauth2"

	"github.com/leor-w/injector"
	"github.com/leor-w/kid/config"
	"github.com/leor-w/kid/utils"

	"golang.org/x/oauth2"
)

type OAuth struct {
	options *Options
	config  *oauth2.Config
}

func (oauth *OAuth) Provide(ctx context.Context) any {
	var confName string
	name, ok := ctx.Value(new(injector.NameKey)).(string)
	if ok && len(name) > 0 {
		confName = "." + name
	}
	confPrefix := fmt.Sprintf("oauth2.google%s", confName)
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

type Option func(o *Options)

const (
	endpointAuth     = "https://open-api.tiktok.com/platform/oauth/connect/"
	endpointToken    = "https://open-api.tiktok.com/oauth/access_token/"
	endpointRefresh  = "https://open-api.tiktok.com/oauth/refresh_token/"
	endpointRevoke   = "https://open-api.tiktok.com/oauth/revoke/"
	endpointUserInfo = "https://open-api.tiktok.com/oauth/userinfo/"
)

func (oauth *OAuth) HandleAuth(code string) (*localOauth2.User, error) {
	token, err := oauth.config.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("授权码换取 token 失败: %s", err.Error())
	}
	client := oauth.config.Client(context.Background(), token)
	// 获取用户信息
	userInfo, err := client.Get(endpointUserInfo)
	if err != nil {
		return nil, fmt.Errorf("获取用户信息失败: %s", err.Error())
	}
	defer userInfo.Body.Close()
	// 解析用户信息
	var user localOauth2.User
	if err := json.NewDecoder(userInfo.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("解析用户信息失败: %s", err.Error())
	}
	return &user, nil
}

func New(opts ...Option) *OAuth {
	options := &Options{}
	for _, opt := range opts {
		opt(options)
	}
	return &OAuth{
		options: options,
		config: &oauth2.Config{
			ClientID:     options.ClientID,
			ClientSecret: options.ClientSecret,
			RedirectURL:  options.RedirectURL,
			Scopes:       options.Scope,
			Endpoint: oauth2.Endpoint{
				AuthURL:   endpointAuth,
				TokenURL:  endpointToken,
				AuthStyle: oauth2.AuthStyleInParams,
			},
		},
	}
}
