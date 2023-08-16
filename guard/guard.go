package guard

import (
	"encoding/json"
	"time"
)

type UserType int

const (
	Admin   UserType = iota + 1 // 后台管理员
	General                     // 普通用户
)

type User struct {
	Uid  int64    `json:"uid"`
	Type UserType `json:"type"`
}

// Guard API 入口守卫
type Guard interface {
	License(*User) (string, error)                // 发行令牌
	GetLicense(UserType, int64) ([]string, error) // 获取用户登录令牌
	Verify(string) (*User, error)                 // 验证令牌
	Cancellation(string) error                    // 吊销令牌
	CancellationAll(UserType, int64) error        // 吊销用户所有令牌
	ExpiresAt(string) int64                       // 获取令牌有效时间
	IssuerAt(string) int64                        // 获取令牌发行时间
}

type TokenInfo struct {
	User
	Token     string `json:"token"`
	ExpiredAt int64  `json:"expired_at"`
	IssuerAt  int64  `json:"issuer_at"`
}

type UserToken struct {
	Tokens []string `json:"tokens"`
}

func (ut *UserToken) Serialize() string {
	utJson, _ := json.Marshal(ut)
	return string(utJson)
}

func (ut *UserToken) Deserialize(utJson string) error {
	if err := json.Unmarshal([]byte(utJson), ut); err != nil {
		return err
	}
	return nil
}

type Store interface {
	GetTokenInfo(license string) (*TokenInfo, error)           // 通过 token 获取 token 的详细信息
	GetUserTokens(uType UserType, uid int64) ([]string, error) // 通过用户类型及用户的 uid 获取用户的 token 列表
	Save(license string, tokenDetail *TokenInfo) error         // 保存用户的 token
	Expired(license string) error                              // 将用户 token 过期
	ExpiredAll(uType UserType, uid int64) error                // 将用户所有 token 过期
	Exist(license string) bool                                 // 检查 token 是否存在 true = 存在 false = 不存在
}

// Blacklist 黑名单列表
type Blacklist interface {
	Black(string, time.Duration) error // 将牌照加入黑名单
	IsBlacked(string) bool             // 检查牌照是否在黑名单列表中
}
