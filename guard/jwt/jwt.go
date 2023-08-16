package jwt

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/leor-w/kid/config"
	"github.com/leor-w/kid/guard"
	"github.com/leor-w/kid/plugin"
	"github.com/leor-w/kid/utils"
)

type Jwt struct {
	blacklist guard.Blacklist `inject:""`
	store     guard.Store     `inject:""`
	options   *Options
}

func (g *Jwt) Provide(ctx context.Context) interface{} {
	var confName string
	name, ok := ctx.Value(plugin.NameKey{}).(string)
	if ok && len(name) > 0 {
		confName += "." + name
	}
	confPrefix := fmt.Sprintf("jwt%s", confName)
	if !config.Exist(confPrefix) {
		panic(fmt.Sprintf("jwt 注入失败: 配置项 [%s] 缺失", confPrefix))
	}
	defaultJwt = New(
		WithIssuer(config.GetString(utils.GetConfigurationItem(confPrefix, "issuer"))),
		WithExpire(config.GetDuration(utils.GetConfigurationItem(confPrefix, "expire"))*time.Hour*24),
		WithSigningMethod(SigningMethod(config.GetString(utils.GetConfigurationItem(confPrefix, "signingMethod")))),
		WithKey([]byte(config.GetString(utils.GetConfigurationItem(confPrefix, "secret")))),
	)
	return defaultJwt
}

type claims struct {
	User *guard.User
	jwt.RegisteredClaims
}

type Option func(*Options)

func (g *Jwt) License(user *guard.User) (string, error) {
	tokenDetail := &guard.TokenInfo{
		User:      *user,
		ExpiredAt: time.Now().Add(g.options.Expire).Unix(),
		IssuerAt:  time.Now().Unix(),
	}
	token := jwt.NewWithClaims(g.getSignMethod(), &claims{
		User: user,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(g.options.Expire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    g.options.Issuer,
		},
	})
	license, err := token.SignedString(g.options.Key)
	if err != nil {
		return "", fmt.Errorf("jwt.License: signed string failed: %w", err)
	}
	tokenDetail.Token = license
	if err := g.store.Save(license, tokenDetail); err != nil {
		return "", err
	}
	return license, nil
}

func (g *Jwt) GetLicense(uType guard.UserType, uid int64) ([]string, error) {
	tokens, err := g.store.GetUserTokens(uType, uid)
	if err != nil {
		return nil, err
	}
	return tokens, nil
}

func (g *Jwt) Verify(license string) (*guard.User, error) {
	token, err := g.checkTokenValidity(license)
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*claims); ok && token.Valid {
		return claims.User, nil
	}
	return nil, fmt.Errorf("jwt.Verify: token valid failed")
}

func (g *Jwt) Cancellation(license string) error {
	token, err := g.checkTokenValidity(license)
	if err != nil {
		return err
	}
	if err := g.store.Expired(license); err != nil {
		return err
	}
	expired := time.Duration(token.Claims.(*claims).ExpiresAt.Unix()-time.Now().Unix()) * time.Second
	if err := g.blacklist.Black(license, expired); err != nil {
		return err
	}
	return nil
}

func (g *Jwt) CancellationAll(uType guard.UserType, uid int64) error {
	return g.store.ExpiredAll(uType, uid)
}

func (g *Jwt) ExpiresAt(license string) int64 {
	token, err := g.checkTokenValidity(license)
	if err != nil {
		return 0
	}
	return token.Claims.(*claims).ExpiresAt.Unix()
}

func (g *Jwt) IssuerAt(license string) int64 {
	token, err := g.checkTokenValidity(license)
	if err != nil {
		return 0
	}
	return token.Claims.(*claims).IssuedAt.Unix()
}

// checkTokenValidity 检查token是否有效
func (g *Jwt) checkTokenValidity(license string) (*jwt.Token, error) {
	// 检查是否过期
	if err := g.isExpired(license); err != nil {
		return nil, err
	}
	// 解析token
	token, err := g.parseToken(license)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (g *Jwt) isExpired(license string) error {
	// 检查是否在黑名单中
	if g.blacklist.IsBlacked(license) {
		return fmt.Errorf("jwt.Verify: token is blacked")
	}
	// 检查是否过期
	if !g.store.Exist(license) {
		return fmt.Errorf("jwt.Verify: token is expired")
	}
	return nil
}

func (g *Jwt) parseToken(license string) (*jwt.Token, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(license, &claims{}, func(token *jwt.Token) (interface{}, error) {
		return g.options.Key, nil
	})
	if err != nil {
		return nil, fmt.Errorf("jwt.Verify: parse claims failed: %w", err)
	}
	if token == nil {
		return nil, fmt.Errorf("jwt.Verify: parse token to claims was nil")
	}
	return token, nil
}

func (g *Jwt) getSignMethod() *jwt.SigningMethodHMAC {
	switch g.options.SigningMethod {
	case SigningMethodSH256:
		return jwt.SigningMethodHS256
	case SigningMethodSH384:
		return jwt.SigningMethodHS384
	case SigningMethodSH512:
		return jwt.SigningMethodHS512
	default:
		return jwt.SigningMethodHS256
	}
}

func New(opts ...Option) *Jwt {
	var options = Options{}
	for _, o := range opts {
		o(&options)
	}
	return &Jwt{
		options: &options,
	}
}

var defaultJwt *Jwt

func Default() *Jwt {
	return defaultJwt
}
