package utils

import (
	"fmt"

	"github.com/nyaruka/phonenumbers"
)

// ValidateInternationalPhoneNumber 验证国际电话号码是否有效
func ValidateInternationalPhoneNumber(number string) (bool, string, error) {
	// 解析电话号码
	parsedNumber, err := phonenumbers.Parse(number, "")
	if err != nil {
		return false, "", fmt.Errorf("错误解析电话号码: %v", err)
	}

	// 验证电话号码是否有效
	if !phonenumbers.IsValidNumber(parsedNumber) {
		return false, "", fmt.Errorf("无效电话号码: %s", number)
	}

	// 格式化为国际标准格式
	formattedNumber := phonenumbers.Format(parsedNumber, phonenumbers.E164)
	return true, formattedNumber, nil
}
