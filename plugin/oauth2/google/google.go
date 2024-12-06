package google

import (
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/api/idtoken"
	"io"
	"net/http"

	"github.com/leor-w/injector"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/leor-w/kid/config"
	plugin "github.com/leor-w/kid/plugin/oauth2"
	"github.com/leor-w/kid/utils"
)

const (
	Profile = "https://www.googleapis.com/auth/userinfo.profile" // 获取用户信息
	Email   = "https://www.googleapis.com/auth/userinfo.email"   // 获取用户邮箱
)

const (
	endpointUserInfo = "https://www.googleapis.com/oauth2/v3/userinfo"
)

type OAuth struct {
	oauthConfig oauth2.Config
	options     Options
}

func (auth *OAuth) Provide(ctx context.Context) any {
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

type Option func(o *Options)

func (auth *OAuth) HandleAuth(code *plugin.VerifyCode) (*plugin.User, error) {
	codeUnescape, err := utils.RecursiveURLDecode(code.Code)
	if err != nil {
		return nil, fmt.Errorf("解码授权码失败: %s", err.Error())
	}
	if code.Code == "" && code.Token != "" {
		// 通过 token 获取用户信息
		ctx := context.Background()
		validator, err := idtoken.NewValidator(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to create validator: %w", err)
		}

		// 验证 ID Token
		payload, err := validator.Validate(ctx, code.Token, "")
		if err != nil {
			return nil, fmt.Errorf("failed to validate ID token: %w", err)
		}
		return &plugin.User{
			UserId:   payload.Subject,
			Email:    payload.Claims["email"].(string),
			EmailVer: payload.Claims["email_verified"].(bool),
			UserName: payload.Claims["name"].(string),
			Picture:  payload.Claims["picture"].(string),
			Locale:   payload.Claims["locale"].(string),
		}, nil
	} else if code.Code != "" {
		// 通过授权码换取 token
		token, err := auth.oauthConfig.Exchange(context.Background(), codeUnescape)
		if err != nil {
			return nil, fmt.Errorf("授权码换取 token 失败: %s", err.Error())
		}
		client := auth.oauthConfig.Client(context.Background(), token)
		// 获取用户信息
		resp, err := client.Get(endpointUserInfo)
		if err != nil {
			return nil, fmt.Errorf("获取用户信息失败: %s", err.Error())
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			bodys, _ := io.ReadAll(resp.Body)
			return nil, fmt.Errorf("获取用户信息失败: %s", string(bodys))
		}
		// 解析用户信息
		var user plugin.User
		if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
			return nil, fmt.Errorf("解析用户信息失败: %s", err.Error())
		}
		return &user, nil
	} else {
		return nil, fmt.Errorf("授权码为空")
	}

}

func (auth *OAuth) GetAuthPageURL() string {
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
