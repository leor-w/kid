package utils

import (
	"crypto/rand"
	"math/big"
)

// RandomSMSCode 随机生成短信验证码
func RandomSMSCode(length int) string {
	if length < 4 {
		length = 4
	}
	const digits = "0123456789"
	code := make([]byte, length)
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return ""
		}
		code[i] = digits[num.Int64()]
	}
	return string(code)
}
