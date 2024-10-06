package third_part

import (
	"context"

	"github.com/markbates/goth/providers/twitterv2"

	"github.com/markbates/goth/providers/tiktok"

	"github.com/markbates/goth/providers/github"

	"github.com/markbates/goth/providers/facebook"

	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/apple"
	"github.com/markbates/goth/providers/google"

	"github.com/leor-w/kid/config"
)

type ThirdPart struct{}

func (tp *ThirdPart) Provide(_ context.Context) any {
	return NewThirdPart()
}

// 配置第三方登录，参考： https://github.com/markbates/goth/blob/master/examples/main.go 这个goth的官方示例，一般只需要传入 key, secret, callback_url 这三个参数即可

func InitProvider() {
	var providers []goth.Provider
	// apple oauth 提供的 scope 有两个，一个是 name，一个是 email
	if config.Exist("third_part.apple") {
		scopeName := apple.ScopeName
		scopeEmail := apple.ScopeEmail
		if config.GetString("third_part.apple.scope_name") != "" {
			scopeName = config.GetString("third_part.apple.scope_name")
		}
		if config.GetString("third_part.apple.scope_email") != "" {
			scopeEmail = config.GetString("third_part.apple.scope_email")
		}
		providers = append(providers, apple.New(
			config.GetString("third_part.apple.key"),
			config.GetString("third_part.apple.secret"),
			config.GetString("third_part.apple.callback_url"),
			nil,
			scopeName,
			scopeEmail,
		))
	}
	// google oauth 提供的 scope 有一个，是 profile
	if config.Exist("third_part.google") {
		providers = append(providers, google.New(
			config.GetString("third_part.google.key"),
			config.GetString("third_part.google.secret"),
			config.GetString("third_part.google.callback_url"),
		))
	}
	if config.Exist("third_part.facebook") {
		providers = append(providers, facebook.New(
			config.GetString("third_part.facebook.key"),
			config.GetString("third_part.facebook.secret"),
			config.GetString("third_part.facebook.callback_url"),
		))
	}
	if config.Exist("third_part.github") {
		providers = append(providers, github.New(
			config.GetString("third_part.github.key"),
			config.GetString("third_part.github.secret"),
			config.GetString("third_part.github.callback_url"),
		))
	}
	if config.Exist("third_part.tiktok") {
		providers = append(providers, tiktok.New(
			config.GetString("third_part.tiktok.key"),
			config.GetString("third_part.tiktok.secret"),
			config.GetString("third_part.tiktok.callback_url"),
		))
	}
	if config.Exist("third_part.twitterv2") {
		providers = append(providers, twitterv2.New(
			config.GetString("third_part.twitterv2.key"),
			config.GetString("third_part.twitterv2.secret"),
			config.GetString("third_part.twitterv2.callback_url"),
		))
	}
	goth.UseProviders(providers...)
}

func NewThirdPart() *ThirdPart {
	InitProvider()
	return &ThirdPart{}
}
