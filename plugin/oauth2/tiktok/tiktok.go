package tiktok

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
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
	options *Options
	config  *oauth2.Config
	rds     *redis.Client `inject:""`
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
		WithClientKey(config.GetString(utils.GetConfigurationItem(confPrefix, "client_key"))),
		WithClientSecret(config.GetString(utils.GetConfigurationItem(confPrefix, "client_secret"))),
		WithRedirectURL(config.GetString(utils.GetConfigurationItem(confPrefix, "redirect_url"))),
		WithScope(config.GetStringSlice(utils.GetConfigurationItem(confPrefix, "scope"))),
	)
}

type Option func(o *Options)

func (oauth *OAuth) HandleOAuth2ByAuthCode(code *plugin.VerifyCode) (*plugin.OAuthUser, error) {
	token, err := oauth.FetchAccessToken(code.Code, code.CodeVerifier)
	if err != nil {
		return nil, err
	}
	// 获取用户信息
	userInfo, err := oauth.GetUserInfo(token.AccessToken, code.Fields)
	if err != nil {
		return nil, err
	}
	if userInfo.Error.Code != "ok" {
		return nil, fmt.Errorf("获取用户信息失败: 错误类型 %s; 错误信息: %s", userInfo.Error.Code, userInfo.Error.Message)
	}
	data, exist := userInfo.Data["user"]
	if !exist {
		return nil, errors.New("获取用户信息失败: 未找到用户信息")
	}
	return &plugin.OAuthUser{
		UserId:   data.OpenId,
		UserName: data.Username,
		Picture:  data.AvatarURL100,
	}, nil
}

// HandleOAuth2ByAPPAuthToken 处理 APP 授权登录
func (oauth *OAuth) HandleOAuth2ByAPPAuthToken(token string) (*plugin.OAuthUser, error) {
	return nil, errors.New("tiktok oauth2: Tiktok 不支持此方法进行验证")
}

func (oauth *OAuth) stateExist(state string) bool {
	return oauth.rds.Exists(GetOAuthURLIdentifierKey(state)).Val() > 0
}

// BuildAuthPageURL 构建授权页面 URL
func (oauth *OAuth) BuildAuthPageURL() (string, error) {
	state := utils.UUID()
	baseUrl, err := url.Parse(endpointAuth)
	if err != nil {
		return "", fmt.Errorf("tiktok oauth2: 解析地址失败: %s", err.Error())
	}
	params := url.Values{}
	params.Add("client_key", oauth.options.ClientKey)
	params.Add("scope", strings.Join(oauth.options.Scope, ","))
	params.Add("response_type", "code")
	params.Add("redirect_uri", oauth.options.RedirectURL)
	params.Add("state", state)
	baseUrl.RawQuery = params.Encode()
	buildUrl := baseUrl.String()
	if err := oauth.rds.Set(GetOAuthURLIdentifierKey(state), buildUrl, time.Minute*30).Err(); err != nil {
		return "", fmt.Errorf("tiktok oauth2: 设置 state 失败: %s", err.Error())
	}
	return buildUrl, nil
}

// executePostRequest 执行 POST 请求的通用函数
func (oauth *OAuth) executePostRequest(endpoint string, data url.Values, result interface{}) error {
	req, err := http.NewRequest("POST", endpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("创建请求失败: %s", err.Error())
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cache-Control", "no-cache")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("发送请求失败: %s", err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("请求失败，状态码：%d", resp.StatusCode)
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("解析响应失败: %s", err.Error())
	}
	return nil
}

func (oauth *OAuth) FetchAccessToken(code, codeVerifier string) (*FetchAccessTokenResp, error) {
	data := url.Values{}
	data.Set("client_key", oauth.options.ClientKey)
	data.Set("client_secret", oauth.options.ClientSecret)
	data.Set("code", code)
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", oauth.options.RedirectURL)
	if len(codeVerifier) > 0 {
		data.Set("code_verifier", codeVerifier)
	}
	var tokenResp FetchAccessTokenResp
	if err := oauth.executePostRequest(endpointToken, data, &tokenResp); err != nil {
		return nil, err
	}
	if len(tokenResp.Error) > 0 {
		return nil, fmt.Errorf("获取 access_token 失败: 错误类型 %s； 错误原因：%s", tokenResp.Error, tokenResp.ErrorDescription)
	}
	return &tokenResp, nil
}

func (oauth *OAuth) RefreshAccessToken(refreshToken string) (*FetchAccessTokenResp, error) {
	data := url.Values{}
	data.Set("client_key", oauth.options.ClientKey)
	data.Set("client_secret", oauth.options.ClientSecret)
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)
	var tokenResp FetchAccessTokenResp
	if err := oauth.executePostRequest(endpointToken, data, &tokenResp); err != nil {
		return nil, err
	}
	if len(tokenResp.Error) > 0 {
		return nil, fmt.Errorf("刷新 access_token 失败: 错误类型 %s； 错误原因：%s", tokenResp.Error, tokenResp.ErrorDescription)
	}
	return &tokenResp, nil
}

func (oauth *OAuth) RevokeToken(accessToken string) error {
	data := url.Values{}
	data.Set("client_key", oauth.options.ClientKey)
	data.Set("client_secret", oauth.options.ClientSecret)
	data.Set("token", accessToken)
	return oauth.executePostRequest(endpointRevoke, data, nil)
}

func (oauth *OAuth) GetUserInfo(accessToken string, fields []string) (*UserInfoResp, error) {
	baseURL, err := url.Parse(endpointUserInfo)
	if err != nil {
		return nil, fmt.Errorf("解析地址失败: %s", err.Error())
	}
	if len(fields) <= 0 {
		fields = []string{"open_id", "union_id", "avatar_url"}
	}
	query := url.Values{}
	query.Set("fields", strings.Join(fields, ","))
	baseURL.RawQuery = query.Encode()
	req, err := http.NewRequest("GET", baseURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %s", err.Error())
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %s", err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			logger.Error("读取响应体失败: %s", err.Error())
		}
		return nil, fmt.Errorf("请求失败，状态码：%d; 错误响应内容：%s", resp.StatusCode, string(body))
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %s", err.Error())
	}
	var userInfoResp UserInfoResp
	if err := json.Unmarshal(body, &userInfoResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %s", err.Error())
	}
	return &userInfoResp, nil
}

// New 创建一个 OAuth 实例
func New(opts ...Option) *OAuth {
	options := &Options{}
	for _, opt := range opts {
		opt(options)
	}
	return &OAuth{
		options: options,
		config: &oauth2.Config{
			ClientID:     options.ClientKey,
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
