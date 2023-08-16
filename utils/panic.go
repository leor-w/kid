package utils

import (
	"fmt"
)

func NoPanic(fn func()) (err error) {
	defer func() {
		if value := recover(); value != nil {
			switch v := value.(type) {
			case error:
				err = v
			default:
				err = fmt.Errorf("%v", v)
			}
		}
	}()
	fn()
	return
}
