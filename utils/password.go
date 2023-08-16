package utils

import (
	"crypto/sha256"
	"encoding/base64"
)

func PasswordEncode(password string) string {
	hasher := sha256.New()
	hasher.Write([]byte(password))
	return base64.StdEncoding.EncodeToString(hasher.Sum(nil))
}

func Verify(plaintext, encrypt string) bool {
	return PasswordEncode(plaintext) == encrypt
}
