package utils

import (
	"fmt"
)

func GetConfigurationItem(prefix, item string) string {
	if len(prefix) == 0 {
		return item
	}
	if len(item) == 0 {
		return prefix
	}
	if prefix[len(prefix)-1] == '.' {
		return prefix + item
	}
	return fmt.Sprintf("%s.%s", prefix, item)
}
