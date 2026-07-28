package constant

import (
	"fmt"
	"strings"

	"github.com/lsjhtang/kid/guard"
)

const (
	tokenDetail     = "auth:token:%s"
	userTokens      = "auth:uid:%d:%d:%s" // auth:uid:<user_type>:<uid>:<token>
	userTokenSearch = "auth:uid:%d:%d"    // auth:uid:<user_type>:<uid>
)

func GetTokenDetailKey(session string) string {
	return fmt.Sprintf(tokenDetail, session)
}

func GetUserTokensKey(uType guard.UserType, uid int64, license string) string {
	return fmt.Sprintf(userTokens, uType, uid, license)
}

func GetUserTokenSearchKey(uType guard.UserType, uid int64) string {
	return fmt.Sprintf(userTokenSearch, uType, uid)
}

func GetTokenByKey(key string) string {
	vals := strings.Split(key, ":")
	return vals[len(vals)-1]
}
