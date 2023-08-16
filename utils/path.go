package utils

import "path"

func GetFileName(file string) string {
	base := path.Base(file)
	return base[:len(base)-len(path.Ext(base))]
}
