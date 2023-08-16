package jwt

import "fmt"

type SigningMethod string

const (
	SigningMethodSH256 SigningMethod = "sha256"
	SigningMethodSH384 SigningMethod = "sha384"
	SigningMethodSH512 SigningMethod = "sha512"
)

const (
	blacklistKey = "auth:blacklist:%s"
)

func GetBlacklistKey(license string) string {
	return fmt.Sprintf(blacklistKey, license)
}
