package token

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/lsjhtang/kid/injector"

	"github.com/lsjhtang/kid/config"
	"github.com/lsjhtang/kid/guard"
	"github.com/lsjhtang/kid/utils"
	"github.com/lsjhtang/kid/utils/signature"
)

type Token struct {
	store   guard.Store `inject:""`
	options *Options
}

func (token *Token) Provide(ctx context.Context) interface{} {
	var confName string
	if name, ok := ctx.Value(injector.NameKey{}).(string); ok && len(name) > 0 {
		confName = "." + name
	}
	confPrefix := fmt.Sprintf("token%s", confName)
	if !config.Exist(confPrefix) {
		panic(fmt.Sprintf("config.yaml file not found configuration item [%s]", confPrefix))
	}
	defaultToken = New(
		WithSecret(config.GetString(utils.GetConfigurationItem(confPrefix, "secret"))),
		WithExpire(config.GetInt(utils.GetConfigurationItem(confPrefix, "expire"))),
	)
	return defaultToken
}

type Option func(*Options)

func (token *Token) License(user *guard.User) (string, string, error) {
	var (
		tokenDetail *guard.TokenInfo
		tokenStr    string
	)
	if user.Uid <= 0 {
		return "", "", errors.New("用户 ID 不能为 0")
	}
	for {
		tokenDetail = &guard.TokenInfo{
			User:      *user,
			ExpiredAt: time.Now().Add(token.options.expire).Unix(),
			IssuerAt:  time.Now().Unix(),
		}
		data, err := signature.EncodeToURLBase64(tokenDetail)
		if err != nil {
			return "", "", err
		}
		tokenStr, err = signature.SignHMACHS384.Sign(data, token.options.secret)
		if err != nil {
			return "", "", err
		}
		tokenDetail.Token = tokenStr
		if token.store.Exist(tokenStr) {
			continue
		}
		break
	}
	if err := token.store.Save(tokenStr, tokenDetail); err != nil {
		return "", "", err
	}
	return tokenStr, "", nil
}

func (token *Token) RefreshLicense(refreshLicense string) (string, error) {
	return "", nil
}

func (token *Token) GetLicense(userType guard.UserType, uid int64) ([]string, error) {
	tokens, err := token.store.GetUserTokens(userType, uid)
	if err != nil {
		return nil, err
	}
	return tokens, nil
}

func (token *Token) Verify(license string) (*guard.User, error) {
	tokenInfo, err := token.checkTokenValidity(license)
	if err != nil {
		return nil, err
	}
	if tokenInfo.ExpiredAt < time.Now().Unix() {
		return nil, fmt.Errorf("token expired")
	}
	return &tokenInfo.User, nil
}

func (token *Token) checkTokenValidity(license string) (*guard.TokenInfo, error) {
	// 检查是否过期
	if !token.store.Exist(license) {
		return nil, fmt.Errorf("jwt.Verify: token is expired")
	}
	// 获取token详情
	tokenInfo, err := token.store.GetTokenInfo(license)
	if err != nil {
		return nil, err
	}
	return tokenInfo, nil
}

func (token *Token) Cancellation(license string) error {
	return token.store.Expired(license)
}

func (token *Token) CancellationAll(uType guard.UserType, uid int64) error {
	return token.store.ExpiredAll(uType, uid)
}

func (token *Token) ExpiresAt(license string) int64 {
	tokenDetail, err := token.store.GetTokenInfo(license)
	if err != nil {
		return 0
	}
	return tokenDetail.ExpiredAt
}

func (token *Token) IssuerAt(license string) int64 {
	tokenDetail, err := token.store.GetTokenInfo(license)
	if err != nil {
		return 0
	}
	return tokenDetail.IssuerAt
}

func New(opts ...Option) *Token {
	options := Options{}
	for _, opt := range opts {
		opt(&options)
	}
	return &Token{
		options: &options,
	}
}

var defaultToken *Token

func Default() *Token {
	return defaultToken
}
