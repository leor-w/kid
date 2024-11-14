package facebook

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/leor-w/injector"
	"github.com/leor-w/kid/config"
	"github.com/leor-w/kid/utils"
	"golang.org/x/oauth2/facebook"

	localOauth2 "github.com/leor-w/kid/plugin/oauth2"
	"golang.org/x/oauth2"
)

type OAuth struct {
	options        Options
	facebookConfig oauth2.Config
}

func (oauth *OAuth) Provide(ctx context.Context) any {
	var confName string
	name, ok := ctx.Value(injector.NameKey{}).(string)
	if ok && len(name) > 0 {
		confName = "." + name
	}
	confPrefix := fmt.Sprintf("oauth2.facebook%s", confName)
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

func (oauth *OAuth) HandleOAuth(code string) (*localOauth2.User, error) {
	token, err := oauth.facebookConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("交换授权码失败: %s", err.Error())
	}
	// 请求 Facebook Graph API 获取用户信息
	userInfoURL := fmt.Sprintf("https://graph.facebook.com/me?fields=id,name,email&access_token=%s", token.AccessToken)
	resp, err := http.Get(userInfoURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("获取用户信息失败: %s", string(body))
	}

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &localOauth2.User{
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
	return &OAuth{
		options: o,
		facebookConfig: oauth2.Config{
			ClientID:     o.ClientID,
			ClientSecret: o.ClientSecret,
			RedirectURL:  o.RedirectURL,
			Scopes:       o.Scope,
			Endpoint:     facebook.Endpoint,
		},
	}
}
