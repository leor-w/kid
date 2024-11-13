package google

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/leor-w/injector"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/leor-w/kid/config"
	localOauth2 "github.com/leor-w/kid/plugin/oauth2"
	"github.com/leor-w/kid/utils"
)

const (
	Profile = "https://www.googleapis.com/auth/userinfo.profile" // 获取用户信息
	Email   = "https://www.googleapis.com/auth/userinfo.email"   // 获取用户邮箱
)

const tokenState = "random"

type OAuth struct {
	oauthConfig oauth2.Config
	options     Options
}

func (auth *OAuth) Provide(ctx context.Context) any {
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

func (auth *OAuth) HandlerAuth(code string) (*localOauth2.User, error) {
	// 通过授权码换取 token
	token, err := auth.oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("授权码换取 token 失败: %s", err.Error())
	}
	client := auth.oauthConfig.Client(context.Background(), token)
	// 获取用户信息
	userInfo, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
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

func (auth *OAuth) GetLoginURL() string {
	return auth.oauthConfig.AuthCodeURL(utils.UUID())
}

func New(opts ...Option) *OAuth {
	o := Options{}
	for _, opt := range opts {
		opt(&o)
	}
	return &OAuth{
		options: o,
		oauthConfig: oauth2.Config{
			ClientID:     o.ClientID,
			ClientSecret: o.ClientSecret,
			RedirectURL:  o.RedirectURL,
			Scopes:       o.Scope,
			Endpoint:     google.Endpoint,
		},
	}
}
