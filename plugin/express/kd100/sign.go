package kd100

import (
	"crypto/md5"
	"fmt"
	"strings"
)

// Sign 数据签名
func (exp *Express) Sign(params string) string {
	hash := md5.Sum([]byte(params))
	return strings.ToUpper(fmt.Sprintf("%x", hash))
}

// Verify 验签
func (exp *Express) Verify(sign, params string) bool {
	return exp.Sign(params+exp.options.Salt) == sign
}
