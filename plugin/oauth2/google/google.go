package google

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/lsjhtang/kid/database/redis"
	"github.com/lsjhtang/kid/injector"
	"github.com/lsjhtang/kid/logger"
	"github.com/spf13/cast"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/idtoken"

	"github.com/lsjhtang/kid/config"
	plugin "github.com/lsjhtang/kid/plugin/oauth2"
	"github.com/lsjhtang/kid/utils"
)

const (
	LockKey = "google.oauth.lock"
)

type OAuth struct {
	oauthConfig oauth2.Config
	options     Options
	rds         *redis.Client `inject:""`
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

// HandleOAuth2ByAuthCode 处理授权码登录授权
func (auth *OAuth) HandleOAuth2ByAuthCode(code *plugin.VerifyCode) (*plugin.OAuthUser, error) {
	if !auth.stateExist(code.State) {
		return nil, fmt.Errorf("google oauth2: 未知的授权来源或授权链接已过期，请重新授权")
	}
	codeUnescape, err := utils.RecursiveURLDecode(code.Code)
	if err != nil {
		return nil, fmt.Errorf("google oauth2: : 解码授权码失败: %s", err.Error())
	}
	// 通过授权码换取 token
	token, err := auth.oauthConfig.Exchange(context.Background(), codeUnescape)
	if err != nil {
		return nil, fmt.Errorf("google oauth2: : 授权码换取 token 失败: %s", err.Error())
	}
	client := auth.oauthConfig.Client(context.Background(), token)
	// 获取用户信息
	resp, err := client.Get(EndpointUserInfo)
	if err != nil {
		return nil, fmt.Errorf("google oauth2: : 获取用户信息失败: %s", err.Error())
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			logger.Errorf("google oauth2: : 关闭获取信息响应流失败: %s", err.Error())
		}
	}(resp.Body)
	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("google oauth2: : 获取用户信息失败: %s", err.Error())
		}
		return nil, fmt.Errorf("google oauth2: : 获取用户信息失败: %s", string(body))
	}
	// 解析用户信息
	var user plugin.OAuthUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("google oauth2: : 解析用户信息失败: %s", err.Error())
	}
	return &user, nil
}

// HandleOAuth2ByAPPAuthToken 处理 APP 授权登录
func (auth *OAuth) HandleOAuth2ByAPPAuthToken(token string) (*plugin.OAuthUser, error) {
	// 通过 token 获取用户信息
	validator, err := idtoken.NewValidator(context.Background())
	if err != nil {
		return nil, fmt.Errorf("google oauth2: : 验证用户信息失败: %w", err)
	}
	// 验证 ID Token
	payload, err := validator.Validate(context.Background(), token, "")
	if err != nil {
		return nil, fmt.Errorf("google oauth2: : 验证用户信息失败: %w", err)
	}
	return &plugin.OAuthUser{
		UserId:   payload.Subject,
		Email:    cast.ToString(payload.Claims["email"]),
		EmailVer: cast.ToBool(payload.Claims["email_verified"]),
		UserName: cast.ToString(payload.Claims["name"]),
		Picture:  cast.ToString(payload.Claims["picture"]),
		Locale:   cast.ToString(payload.Claims["locale"]),
	}, nil
}

// stateExist 判断 state 是否存在
func (auth *OAuth) stateExist(state string) bool {
	return auth.rds.Exists(GetOAuthURLIdentifierKey(state)).Val() > 0
}

// BuildAuthPageURL 构建授权页面 URL
func (auth *OAuth) BuildAuthPageURL() (string, error) {
	state := utils.UUID()
	buildUrl := auth.oauthConfig.AuthCodeURL(state)
	if err := auth.rds.Set(GetOAuthURLIdentifierKey(state), buildUrl, time.Minute*30).Err(); err != nil {
		return "", fmt.Errorf("google oauth2: : 保存授权链接映射关系报错: %s", err.Error())
	}
	return buildUrl, nil
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
