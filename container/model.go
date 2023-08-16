package container

import (
	"strings"
)

// InjectTag 注入的所有 tag
type InjectTag map[string]string

// parseTag 解析 tag
func parseTag(tag string) InjectTag {
	it := make(InjectTag)
	if tag == "" {
		return nil
	}
	if strings.Contains(tag, ",") {
		tags := strings.Split(tag, ",")
		for _, t := range tags {
			if t == "" || !strings.Contains(t, ":") {
				continue
			}
			k, v := parseValue(t)
			it[k] = v
		}
	} else {
		k, v := parseValue(tag)
		it[k] = v
	}
	return it
}

func parseValue(t string) (key, value string) {
	items := strings.Split(t, ":")
	if len(items) == 2 {
		key = items[0]
		value = items[1]
	}
	return
}

// getAlias 获取注入的别名
func (t InjectTag) getAlias() string {
	return t[tagAlias]
}

// getOpts 获取注入的选项
func (t InjectTag) getOpts() string {
	return t[tagOpts]
}

// hasAlias 是否有别名
func (t InjectTag) hasAlias() bool {
	return t[tagAlias] != ""
}

// hasOpts 是否有选项
func (t InjectTag) hasOpts() bool {
	return t[tagOpts] != ""
}

// hasTag 是否有 tag
func (t InjectTag) hasTag() bool {
	return t.hasAlias() || t.hasOpts()
}

// options 获取选项 slice
func (t InjectTag) options() []string {
	return strings.Split(t.getOpts(), "|")
}

// containOption 是否包含选项
func (t InjectTag) containOption(opt string) bool {
	if !t.hasOpts() {
		return false
	}
	if strings.Contains(t.getOpts(), opt) {
		return true
	}
	return false
}
