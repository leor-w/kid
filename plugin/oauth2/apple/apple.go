package apple

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lsjhtang/kid/injector"
	"github.com/spf13/cast"
	"golang.org/x/oauth2"

	"github.com/lsjhtang/kid/config"
	"github.com/lsjhtang/kid/database/redis"
	plugin "github.com/lsjhtang/kid/plugin/oauth2"
	"github.com/lsjhtang/kid/utils"
)

type OAuth struct {
	oauthConfig *oauth2.Config
	options     *Options
	secret      string
	rds         *redis.Client `inject:""`
}

type Option func(o *Options)

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
		WithKeyId(config.GetString(utils.GetConfigurationItem(confPrefix, "key_id"))),
		WithKeySecret(config.GetString(utils.GetConfigurationItem(confPrefix, "key_secret"))),
		WithKeySecretFile(config.GetString(utils.GetConfigurationItem(confPrefix, "key_secret_file"))),
		WithTeamId(config.GetString(utils.GetConfigurationItem(confPrefix, "team_id"))),
		WithRedirectURL(config.GetString(utils.GetConfigurationItem(confPrefix, "redirect_url"))),
		WithScope(config.GetStringSlice(utils.GetConfigurationItem(confPrefix, "scope"))...),
	)
}

func (oauth *OAuth) HandleOAuth2ByAuthCode(code *plugin.VerifyCode) (*plugin.OAuthUser, error) {
	token, err := oauth.oauthConfig.Exchange(context.Background(), code.Code)
	if err != nil {
		return nil, fmt.Errorf("apple oauth2: 通过授权码换取 token 失败: %s", err.Error())
	}
	idToken := cast.ToString(token.Extra("id_token"))
	if idToken == "" {
		return nil, errors.New("apple oauth2: 未找到 id_token")
	}
	return oauth.parseIdToken(idToken)
}

func (oauth *OAuth) HandleOAuth2ByAPPAuthToken(idToken string) (*plugin.OAuthUser, error) {
	return oauth.parseIdToken(idToken)
}

func (oauth *OAuth) BuildAuthPageURL() (string, error) {
	state := utils.UUID()
	baseUrl, err := url.Parse(BuildOAuthPageURL)
	if err != nil {
		return "", fmt.Errorf("apple oauth2: 解析 URL 失败: %s", err.Error())
	}
	params := url.Values{}
	params.Add("response_type", "code")
	params.Add("response_mode", "form_post")
	params.Add("client_id", oauth.options.ClientId)
	params.Add("redirect_uri", oauth.options.RedirectURL)
	params.Add("state", state)
	params.Add("scope", strings.Join(oauth.options.Scopes, " "))
	baseUrl.RawQuery = params.Encode()
	buildUrl := baseUrl.String()
	if err := oauth.rds.Set(GetOAuthURLIdentifierKey(state), buildUrl, time.Minute*30).Err(); err != nil {
		return "", fmt.Errorf("apple oauth2: 设置 state 失败: %s", err.Error())
	}
	return buildUrl, nil
}

func (oauth *OAuth) parseIdToken(idToken string) (*plugin.OAuthUser, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(idToken, jwt.MapClaims{})
	if err != nil {
		return nil, fmt.Errorf("apple OAuth: Token 解析失败: %s", err.Error())
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return &plugin.OAuthUser{
			UserId:   cast.ToString(claims["sub"]),
			Email:    cast.ToString(claims["email"]),
			EmailVer: cast.ToBool(claims["email_verified"]),
			UserName: cast.ToString(claims["name"]),
			Locale:   cast.ToString(claims["locale"]),
		}, nil
	}
	return nil, errors.New("apple OAuth2: Token 验证失败")
}

// 生成 client_secret
func (oauth *OAuth) generateClientSecret() string {
	var (
		privateKey []byte
		err        error
	)
	if oauth.options.KeySecret != "" {
		privateKey = []byte(oauth.options.KeySecret)
	} else {
		privateKey, err = os.ReadFile(oauth.options.KeySecretFile)
		if err != nil {
			panic(fmt.Sprintf("Failed to read private key file: %v", err))
		}
	}

	// 创建 JWT Token
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"iss": oauth.options.TeamId,             // Team ID
		"iat": time.Now().Unix(),                // Issued At
		"exp": time.Now().Add(time.Hour).Unix(), // Expiration
		"aud": "https://appleid.apple.com",      // Audience
		"sub": oauth.options.ClientId,           // Client ID
	})

	// 添加 Key ID 到 Header
	token.Header["kid"] = oauth.options.KeyId

	// 使用私钥签名
	privateKeyParsed, err := jwt.ParseECPrivateKeyFromPEM(privateKey)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse private key: %v", err))
	}

	clientSecret, err := token.SignedString(privateKeyParsed)
	if err != nil {
		panic(fmt.Sprintf("Failed to sign client secret: %v", err))
	}

	return clientSecret
}

func (oauth *OAuth) init() error {
	if oauth.options.KeySecret == "" && oauth.options.KeySecretFile == "" {
		return errors.New("Apple OAuth2: 客户端密钥不能为空")
	}
	if oauth.options.ClientId == "" {
		return errors.New("Apple OAuth2: 客户端ID不能为空")
	}
	if oauth.options.KeyId == "" {
		return errors.New("Apple OAuth2: 密钥ID不能为空")
	}
	if oauth.options.TeamId == "" {
		return errors.New("Apple OAuth2: 团队ID不能为空")
	}
	if oauth.options.RedirectURL == "" {
		return errors.New("Apple OAuth2: 重定向URL不能为空")
	}
	oauth.secret = oauth.generateClientSecret()
	return nil
}

func New(opts ...Option) *OAuth {
	var o Options
	for _, opt := range opts {
		opt(&o)
	}
	var oauth = &OAuth{
		options: &o,
	}
	if err := oauth.init(); err != nil {
		panic(fmt.Sprintf("Apple OAuth2: 初始化失败: %s", err.Error()))
	}
	oauth.oauthConfig = &oauth2.Config{
		ClientID:     o.ClientId,
		ClientSecret: oauth.secret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  BuildOAuthPageURL,
			TokenURL: EndpointTokenURL,
		},
		RedirectURL: o.RedirectURL,
		Scopes:      o.Scopes,
	}
	return oauth
}
