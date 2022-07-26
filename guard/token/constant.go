package token

import (
	"fmt"
	"github.com/leor-w/kid/guard"
)

const (
	Key   = "auth.token.%s"
	IDKey = "auth.%d.%d" // auth.type.id
)

func GetTokenKey(session string) string {
	return fmt.Sprintf(Key, session)
}

func GetTokenIdKey(userType guard.UserType, id int64) string {
	return fmt.Sprintf(IDKey, userType, id)
}
