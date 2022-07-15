package jwt

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/leor-w/kid/config"
	"github.com/leor-w/kid/guard"
	"time"
)

type Jwt struct {
	Blacklist guard.Blacklist
	options   *Options
}

func (guard *Jwt) Provide() interface{} {
	if !config.Exist("jwt") {
		panic("not found [jwt] in config")
	}
	return New(
		WithIssuer(config.GetString("jwt.issuer")),
		WithExpire(config.GetDuration("jwt.expire")*time.Hour*24),
		WithSigningMethod(SigningMethod(config.GetString("jwt.signingMethod"))),
		WithKey([]byte(config.GetString("jwt.secret"))),
	)
}

type kidClaims struct {
	User *guard.User
	jwt.StandardClaims
}

type Option func(*Options)

func (guard *Jwt) License(user *guard.User) (string, error) {
	token := jwt.NewWithClaims(guard.getSignMethod(), &kidClaims{
		User: user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(guard.options.Expire).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    guard.options.Issuer,
		},
	})
	var err error
	license, err := token.SignedString(guard.options.SigningMethod)
	if err != nil {
		return "", fmt.Errorf("jwt.License: signed string failed: %w", err)
	}
	return license, nil
}

func (guard *Jwt) Verify(license string) (*guard.User, error) {
	token, err := guard.parseToken(license)
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*kidClaims); ok && token.Valid {
		return claims.User, nil
	}
	return nil, fmt.Errorf("jwt.Verify: token valid failed")
}

func (guard *Jwt) Cancellation(license string) error {
	token, err := guard.parseToken(license)
	if err != nil {
		return err
	}
	return guard.Blacklist.Black(license, time.Duration(token.Claims.(*kidClaims).ExpiresAt-time.Now().Unix()))
}

func (guard *Jwt) ExpiresAt(license string) int64 {
	token, err := guard.parseToken(license)
	if err != nil {
		return 0
	}
	return token.Claims.(*kidClaims).ExpiresAt
}

func (guard *Jwt) IssuerAt(license string) int64 {
	token, err := guard.parseToken(license)
	if err != nil {
		return 0
	}
	return token.Claims.(*kidClaims).IssuedAt
}

func (guard *Jwt) parseToken(license string) (*jwt.Token, error) {
	isBlack, err := guard.Blacklist.Checklist(license)
	if err != nil {
		return nil, fmt.Errorf("jwt.Verify: check blacklist failed: %w", err)
	}
	if isBlack {
		return nil, errors.New("jwt.Verify: token is invalid")
	}
	token, err := jwt.ParseWithClaims(license, &kidClaims{}, func(token *jwt.Token) (interface{}, error) {
		return guard.options.Key, nil
	})
	if err != nil {
		return nil, fmt.Errorf("jwt.Verify: parse claims failed: %w", err)
	}
	if token == nil {
		return nil, fmt.Errorf("jwt.Verify: parse token to claims was nil")
	}
	return token, nil
}

func (guard *Jwt) getSignMethod() *jwt.SigningMethodHMAC {
	switch guard.options.SigningMethod {
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
