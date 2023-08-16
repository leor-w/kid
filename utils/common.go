package utils

import (
	"fmt"
	"os"
	"strings"
)

func ContainString(src []string, tag string) bool {
	if len(src) == 0 {
		return false
	}
	for i := range src {
		if tag == src[i] {
			return true
		}
	}
	return false
}

func ContainInt(src []int, tag int) bool {
	if len(src) == 0 {
		return false
	}
	for i := range src {
		if src[i] == tag {
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

func ParseConfig(conf string) ([]string, error) {
	split := strings.Split(conf, ",")
	if len(split) > 0 {
		for i := range split {
			split[i] = strings.Trim(split[i], " \n\t")
		}
		return split, nil
	}
	return nil, fmt.Errorf("非法的配置文件路径")
}

func FileExist(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if err != nil {
		return false, err
	}
	if os.IsNotExist(err) {
		return false, fmt.Errorf("文件不存在")
	}
	return true, nil
}
