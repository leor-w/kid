package utils

func ContainString(src []string, tag string) bool {
	for i := range src {
		if tag == src[i] {
			return true
		}
	}
	return false
}
