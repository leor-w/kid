package utils

import "strings"

func ContainString(src []string, tag string) bool {
	for i := range src {
		if tag == src[i] {
			return true
		}
	}
	return false
}

func CamelToSnake(camel string) string {
	data := make([]byte, 0, len(camel)*2)
	j := false
	num := len(camel)
	for i := 0; i < num; i++ {
		d := camel[i]
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	return strings.ToLower(string(data[:]))
}

func SnakeToCamel(snake string) string {
	data := make([]byte, 0, len(snake))
	j := false
	k := false
	num := len(snake) - 1
	for i := 0; i < num; i++ {
		d := snake[i]
		if !k && d >= 'A' && d <= 'Z' {
			k = true
		}
		if d >= 'a' && d <= 'z' && (j || !k) {
			d = d - 32
			j = false
			k = true
		}
		if k && d == '_' && num > i && snake[i+1] >= 'a' && snake[i+1] <= 'z' {
			j = true
			continue
		}
		data = append(data, d)
	}
	return string(data[:])
}
