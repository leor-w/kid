package utils

import (
	"fmt"
	"strings"

	idvalidator "github.com/guanguans/id-validator"
)

// ValidateChinaIdCard 验证中国身份证号码（可以验证包括15位、18位、末尾为x的中国身份证号码，港澳居民证，台湾居民证等）
func ValidateChinaIdCard(idNum string, strict bool) bool {
	return idvalidator.IsValid(idNum, strict)
}

// MaskIDCard 对身份证号码进行脱敏处理
func MaskIDCard(id string) (string, error) {
	length := len(id)

	// 检查身份证长度是否合法
	if length != 15 && length != 18 {
		return "", fmt.Errorf("身份证号码长度非法")
	}

	// 脱敏处理
	prefix := id[:3]
	suffix := id[length-4:]
	masked := strings.Repeat("*", length-7)

	return prefix + masked + suffix, nil
}
