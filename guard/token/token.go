package token

import (
	"fmt"
	"github.com/leor-w/kid/config"
	"github.com/leor-w/kid/guard"
	"github.com/leor-w/kid/utils/signature"
	"time"
)

type Token struct {
	Store   Store `inject:""`
	options *Options
}

func (guard *Token) Provide() interface{} {
	if !config.Exist("token") {
		panic("not found [token] in config")
	}
	token = New(
		WithSecret(config.GetString("token.secret")),
		WithExpire(config.GetInt("token.expire")),
	)
	return token
}

type Option func(*Options)

type tokenInfo struct {
	*guard.User
	Token    string `json:"token"`
	ExpireAt int64  `json:"expire_at"`
	IssuerAt int64  `json:"issuer_at"`
}

func (guard *Token) License(user *guard.User) (string, error) {
	var (
		info     *tokenInfo
		tokenStr string
	)
	for {
		info = &tokenInfo{
			User:     user,
			ExpireAt: time.Now().Add(guard.options.expire).Unix(),
			IssuerAt: time.Now().Unix(),
		}
		data, err := signature.EncodeToURLBase64(info)
		if err != nil {
			return "", err
		}
		tokenStr, err = signature.SignHMACHS384.Sign(data, guard.options.secret)
		if err != nil {
			return "", err
		}
		info.Token = tokenStr
		if guard.Store.Exist(tokenStr) {
			continue
		}
		break
	}
	if err := guard.Store.Save(tokenStr, info); err != nil {
		return "", err
	}
	return tokenStr, nil
}

func (guard *Token) Verify(license string) (*guard.User, error) {
	token, err := guard.Store.Get(license)
	if err != nil {
		return nil, err
	}
	if token.ExpireAt < time.Now().Unix() {
		return nil, fmt.Errorf("token expired")
	}
	return token.User, nil
}

func (guard *Token) Cancellation(license string) error {
	return guard.Store.Expired(license)
}

func (guard *Token) ExpiresAt(license string) int64 {
	token, err := guard.Store.Get(license)
	if err != nil {
		return 0
	}
	return token.ExpireAt
}

func (guard *Token) IssuerAt(license string) int64 {
	token, err := guard.Store.Get(license)
	if err != nil {
		return 0
	}
	return token.IssuerAt
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

var token *Token

func Default() *Token {
	return token
}
