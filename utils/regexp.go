package utils

import "regexp"

// RegexpMatchPhone 正则匹配手机号
func RegexpMatchPhone(phone string) bool {
	regRuler := "^1[3456789]{1}\\d{9}$"
	reg := regexp.MustCompile(regRuler)
	return reg.MatchString(phone)
}
