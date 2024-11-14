package apple

import (
	"context"
	"fmt"
	"os"

	"github.com/leor-w/injector"

	"github.com/Timothylock/go-signin-with-apple/apple"

	"github.com/leor-w/kid/config"
	localOauth2 "github.com/leor-w/kid/plugin/oauth2"
	"github.com/leor-w/kid/utils"
)

type OAuth struct {
	client  *apple.Client
	options *Options
}

type Option func(o *Options)

func (oauth *OAuth) Provide(ctx context.Context) any {
	var confName string
	name, ok := ctx.Value(injector.NameKey{}).(string)
	if ok && len(name) > 0 {
		confName = "." + name
	}
	confPrefix := fmt.Sprintf("oauth2.apple%s", confName)
	if !config.Exist(confPrefix) {
		panic(fmt.Sprintf("配置文件为找到 [%s.*]，请检查配置文件", confPrefix))
	}
	return New(
		WithClientID(config.GetString(utils.GetConfigurationItem(confPrefix, "client_id"))),
		WithKeyId(config.GetString(utils.GetConfigurationItem(confPrefix, "client_key"))),
		WithClientSecret(config.GetString(utils.GetConfigurationItem(confPrefix, "client_secret"))),
		WithClientSecretFile(config.GetString(utils.GetConfigurationItem(confPrefix, "client_secret_file"))),
		WithTeamId(config.GetString(utils.GetConfigurationItem(confPrefix, "team_id"))),
		WithRedirectURL(config.GetString(utils.GetConfigurationItem(confPrefix, "redirect_url"))),
	)
}

func (oauth *OAuth) HandlerAuth(code string) (*localOauth2.User, error) {
	req := apple.AppValidationTokenRequest{
		ClientID:     oauth.options.ClientId,
		ClientSecret: oauth.options.ClientSecret,
		Code:         code,
	}
	var resp apple.ValidationResponse
	if err := oauth.client.VerifyAppToken(context.Background(), req, &resp); err != nil {
		return nil, fmt.Errorf("验证授权码失败: %s", err.Error())
	}
	if resp.Error != "" {
		return nil, fmt.Errorf("验证授权码失败: %s", resp.Error)
	}
	unique, err := apple.GetUniqueID(resp.IDToken)
	if err != nil {
		return nil, fmt.Errorf("获取唯一标识失败: %s", err.Error())
	}
	claims, err := apple.GetClaims(resp.IDToken)
	if err != nil {
		return nil, fmt.Errorf("获取用户信息失败: %s", err.Error())
	}
	return &localOauth2.User{
		UserId:   unique,
		Email:    (*claims)["email"].(string),
		EmailVer: (*claims)["email_verified"].(bool),
		UserName: (*claims)["name"].(string),
		Locale:   (*claims)["locale"].(string),
	}, nil
}

func New(opts ...Option) *OAuth {
	var o Options
	for _, opt := range opts {
		opt(&o)
	}
	if o.ClientSecret == "" && o.ClientSecretFile == "" {
		panic("Apple OAuth2: 客户端密钥不能为空")
	}
	if o.ClientSecretFile != "" {
		secretBytes, err := os.ReadFile(o.ClientSecretFile)
		if err != nil {
			panic(fmt.Sprintf("读取客户端密钥文件失败: %s", err.Error()))
		}
		secret, err := apple.GenerateClientSecret(string(secretBytes), o.TeamId, o.ClientId, o.KeyId)
		if err != nil {
			panic(fmt.Sprintf("生成客户端密钥失败: %s", err.Error()))
		}
		o.ClientSecret = secret
	}
	if o.ClientId == "" {
		panic("Apple OAuth2: 客户端ID不能为空")
	}
	if o.KeyId == "" {
		panic("Apple OAuth2: 密钥ID不能为空")
	}
	if o.TeamId == "" {
		panic("Apple OAuth2: 团队ID不能为空")
	}
	if o.RedirectURL == "" {
		panic("Apple OAuth2: 重定向URL不能为空")
	}
	var oauth = &OAuth{
		client:  apple.New(),
		options: &o,
	}
	return oauth
}
