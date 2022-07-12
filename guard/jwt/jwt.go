package jwt

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/leor-w/kid"
	"time"
)

type Jwt struct {
	license   string
	claims    *kidClaims
	Blacklist kid.Blacklist
	options   *Options
}

type kidClaims struct {
	User *kid.User
	jwt.StandardClaims
}

type Option func(*Options)

func (guard *Jwt) License(user *kid.User) (string, error) {
	token := jwt.NewWithClaims(guard.getSignMethod(), &kidClaims{
		User: user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(guard.options.Expire).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    guard.options.Issuer,
		},
	})
	var err error
	guard.license, err = token.SignedString(guard.options.SigningMethod)
	if err != nil {
		return "", fmt.Errorf("jwt.License: signed string failed: %w", err)
	}
	return guard.license, nil
}

func (guard *Jwt) Verify(license string) (*kid.User, error) {
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
	if claims, ok := token.Claims.(*kidClaims); ok && token.Valid {
		guard.claims = claims
		return claims.User, nil
	}
	return nil, fmt.Errorf("jwt.Verify: token valid failed")
}

func (guard *Jwt) Cancellation() error {
	return guard.Blacklist.Black(guard.license, time.Duration(guard.claims.ExpiresAt-time.Now().Unix()))
}

func (guard *Jwt) ExpiresAt() int64 {
	return guard.claims.ExpiresAt
}

func (guard *Jwt) IssuerAt() int64 {
	return guard.claims.IssuedAt
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
