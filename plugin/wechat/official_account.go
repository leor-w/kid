package wechat

import (
	"fmt"
	"time"

	"github.com/silenceper/wechat/v2/officialaccount/oauth"

	"github.com/leor-w/kid/utils"
)

// BuildOauthURL 构建微信授权URL
func (w *Wechat) BuildOauthURL(redirectURL, scope string) (string, error) {
	state := utils.UUID()
	if w.cache != nil {
		if err := w.cache.Set(getWechatOauthStateKey(state), state, time.Minute*30); err != nil {
			return "", err
		}
	}
	return w.official.GetOauth().GetRedirectURL(redirectURL, scope, state)
}

// GetOAuthUserInfo 获取微信用户信息
func (w *Wechat) GetOAuthUserInfo(code string) (*oauth.UserInfo, error) {
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
