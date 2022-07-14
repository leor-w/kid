package token

import "fmt"

const (
	Key   = "auth.token.%s"
	IDKey = "auth.id.%d"
)

func GetTokenKey(session string) string {
	return fmt.Sprintf(Key, session)
}

func GetTokenIdKey(id int64) string {
	return fmt.Sprintf(IDKey, id)
}
