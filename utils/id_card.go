package utils

import idvalidator "github.com/guanguans/id-validator"

// ValidateChinaIdCard 验证中国身份证号码（可以验证包括15位、18位、末尾为x的中国身份证号码，港澳居民证，台湾居民证等）
func ValidateChinaIdCard(idNum string, strict bool) bool {
	return idvalidator.IsValid(idNum, strict)
}
