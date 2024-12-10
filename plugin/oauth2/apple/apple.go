package apple

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/spf13/cast"

	"github.com/golang-jwt/jwt/v5"

	"github.com/leor-w/kid/logger"

	"github.com/Timothylock/go-signin-with-apple/apple"
	"github.com/leor-w/injector"

	"github.com/leor-w/kid/config"
	"github.com/leor-w/kid/database/redis"
	plugin "github.com/leor-w/kid/plugin/oauth2"
	"github.com/leor-w/kid/utils"
)

type OAuth struct {
	client  *apple.Client
	options *Options
	rds     *redis.Client `inject:""`
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
		WithKeyId(config.GetString(utils.GetConfigurationItem(confPrefix, "client_key"))),
		WithClientSecret(config.GetString(utils.GetConfigurationItem(confPrefix, "client_secret"))),
		WithClientSecretFile(config.GetString(utils.GetConfigurationItem(confPrefix, "client_secret_file"))),
		WithTeamId(config.GetString(utils.GetConfigurationItem(confPrefix, "team_id"))),
		WithRedirectURL(config.GetString(utils.GetConfigurationItem(confPrefix, "redirect_url"))),
	)
}

func (oauth *OAuth) HandlerAuth(code string) (*plugin.OAuthUser, error) {
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
	return &plugin.OAuthUser{
		UserId:   unique,
		Email:    (*claims)["email"].(string),
		EmailVer: (*claims)["email_verified"].(bool),
		UserName: (*claims)["name"].(string),
		Locale:   (*claims)["locale"].(string),
	}, nil
}

func (oauth *OAuth) HandleOAuth2ByAuthCode(code *plugin.VerifyCode) (*plugin.OAuthUser, error) {
	publicKey, err := oauth.getPublicKeyByAccessToken(code.IdToken)
	if err != nil {
		return nil, fmt.Errorf("apple oauth2: 获取公钥失败: %s", err.Error())
	}
	claims, err := VerifyIDToken(code.IdToken, publicKey)
	if err != nil {
		return nil, fmt.Errorf("apple oauth2: 验证 id_token 失败: %s", err.Error())
	}
	return &plugin.OAuthUser{
		UserId: cast.ToString(claims["sub"]),
		Email:  cast.ToString(claims["email"]),
	}, nil
}

func (oauth *OAuth) HandleOAuth2ByAPPAuthToken(token string) (*plugin.OAuthUser, error) {
	return oauth.HandlerAuth(token)
}

func (oauth *OAuth) BuildAuthPageURL() (string, error) {
	state := utils.UUID()
	baseUrl, err := url.Parse(BuildOAuthPageURL)
	if err != nil {
		return "", fmt.Errorf("apple oauth2: 解析 URL 失败: %s", err.Error())
	}
	params := url.Values{}
	params.Add("client_id", oauth.options.ClientId)
	params.Add("redirect_uri", url.QueryEscape(oauth.options.RedirectURL))
	params.Add("state", state)
	baseUrl.RawQuery = params.Encode()
	buildUrl := baseUrl.String()
	if err := oauth.rds.Set(GetOAuthURLIdentifierKey(state), buildUrl, time.Minute*30).Err(); err != nil {
		return "", fmt.Errorf("apple oauth2: 设置 state 失败: %s", err.Error())
	}
	return buildUrl, nil
}

func (oauth *OAuth) getPublicKeyByAccessToken(idToken string) ([]byte, error) {
	resp, err := http.Get(EndpointPublicKeyURL)
	if err != nil {
		return nil, fmt.Errorf("apple oauth2: 获取公钥失败: %s", err.Error())
	}
	defer func(resp *http.Response) {
		if err := resp.Body.Close(); err != nil {
			logger.Errorf("apple oauth2: 关闭响应流失败: %s", err.Error())
		}
	}(resp)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("apple oauth2: 获取公钥失败: %s", resp.Status)
	}
	var keysResp GetPublicKeyResp
	if err := json.NewDecoder(resp.Body).Decode(&keysResp); err != nil {
		return nil, fmt.Errorf("apple oauth2: 解析公钥失败: %s", err.Error())
	}
	kid := ExtractKidFromIDToken(idToken)
	for _, key := range keysResp.Keys {
		if key.Kid == kid {
			return ConvertKeyToPEM(key.N, key.E)
		}
	}
	return nil, fmt.Errorf("apple oauth2: 未找到指定的公钥")
}

// ExtractKidFromIDToken 从 id_token 提取 kid
func ExtractKidFromIDToken(idToken string) string {
	parts := strings.Split(idToken, ".")
	if len(parts) != 3 {
		return ""
	}
	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return ""
	}
	var header struct {
		Kid string `json:"kid"`
	}
	err = json.Unmarshal(headerBytes, &header)
	if err != nil {
		return ""
	}
	return header.Kid
}

// ConvertKeyToPEM 将 Apple 公钥的 n 和 e 转换为 PEM 格式
func ConvertKeyToPEM(nStr, eStr string) ([]byte, error) {
	// 解码 base64 编码的 n 和 e
	nBytes, err := base64.RawURLEncoding.DecodeString(nStr)
	if err != nil {
		return nil, err
	}
	eBytes, err := base64.RawURLEncoding.DecodeString(eStr)
	if err != nil {
		return nil, err
	}

	// 转换 eBytes 为整数
	var eInt int
	for _, b := range eBytes {
		eInt = eInt<<8 + int(b)
	}

	// 创建 RSA 公钥
	pubKey := &rsa.PublicKey{
		N: new(big.Int).SetBytes(nBytes),
		E: eInt,
	}

	// 转换为 PEM 格式
	pubASN1, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		return nil, err
	}
	pubPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubASN1,
	})

	return pubPEM, nil
}

// VerifyIDToken 验证 id_token
func VerifyIDToken(idToken string, publicKeyPEM []byte) (map[string]interface{}, error) {
	// 解析 PEM 公钥
	block, _ := pem.Decode(publicKeyPEM)
	if block == nil {
		return nil, errors.New("failed to parse PEM block")
	}
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	// 解析和验证 JWT
	token, err := jwt.Parse(idToken, func(token *jwt.Token) (interface{}, error) {
		// 确保签名方法是 RSA
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return pubKey.(*rsa.PublicKey), nil
	})
	if err != nil {
		return nil, err
	}

	// 检查 token 是否有效
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
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
