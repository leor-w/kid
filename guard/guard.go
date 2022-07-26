package guard

import "time"

type UserType int

const (
	Admin UserType = iota + 1
	GeneralUser
)

type User struct {
	Id   int64    `json:"id"`
	Type UserType `json:"type"`
}

// Guard API 入口守卫
type Guard interface {
	License(*User) (string, error) // 发行牌照
	Verify(string) (*User, error)  // 验证牌照
	Cancellation(string) error     // 吊销牌照
	ExpiresAt(string) int64        // 获取牌照有效时间
	IssuerAt(string) int64         // 获取牌照发行时间
}

// Blacklist 黑名单列表
type Blacklist interface {
	Black(string, time.Duration) error // 将牌照加入黑名单
	Checklist(string) (bool, error)    // 检查牌照是否在黑名单列表中
}
