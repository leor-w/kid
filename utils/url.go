package utils

import (
	"fmt"
	"net/url"
)

// RecursiveURLDecode 尝试递归解码 URL 编码的字符串，直到解码后的字符串与原字符串相同，返回解码结果。如果原始值无法解码，返回错误。
func RecursiveURLDecode(encodedStr string) (string, error) {
	decodedStr, err := url.QueryUnescape(encodedStr)
	if err != nil {
		return "", fmt.Errorf("解码失败: %v", err)
	}
	// 如果解码后的字符串与原字符串不同，继续递归解码
	if decodedStr != encodedStr {
		return RecursiveURLDecode(decodedStr)
	}
	// 如果解码后的字符串与原字符串相同，返回解码结果
	return decodedStr, nil
}
