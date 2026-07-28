package wechat

import (
	"fmt"
	"time"

	"github.com/silenceper/wechat/v2/officialaccount/oauth"
	"github.com/spf13/cast"

	"github.com/lsjhtang/kid/utils"
)

// BuildOauthURL 构建微信授权URL
func (w *Wechat) BuildOauthURL(redirectURL, scope string) (string, error) {
	state := utils.UUID()
	if w.cache != nil {
		if err := w.cache.Set(getWechatOauthStateKey(state), state, time.Minute*30); err != nil {
			return "", err
		}
	}
	return w.OfficialAccount().GetOauth().GetRedirectURL(redirectURL, scope, state)
}

// GetOAuthUserInfo 获取微信用户信息
func (w *Wechat) GetOAuthUserInfo(code, state string) (*oauth.UserInfo, error) {
	if w.cache != nil {
		key := getWechatOauthStateKey(state)
		rest := w.cache.Get(key)
		if rest == nil {
			return nil, fmt.Errorf("微信公众号 获取用户信息错误: %w", fmt.Errorf("state 无效"))
		}
		if cast.ToString(rest) != state {
			return nil, fmt.Errorf("微信公众号 获取用户信息错误: %w", fmt.Errorf("state 无效"))
		}
	}
	res, err := w.OfficialAccount().GetOauth().GetUserAccessToken(code)
	if err != nil {
		return nil, err
	}
	userInfo, err := w.OfficialAccount().GetOauth().GetUserInfo(res.AccessToken, res.OpenID, "en")
	if err != nil {
		return nil, fmt.Errorf("微信公众号 获取用户信息错误: %w", err)
	}
	return &userInfo, nil
}
