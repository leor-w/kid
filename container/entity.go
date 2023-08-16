package container

import (
	"reflect"
	"sync"

	"github.com/leor-w/utils"
)

// entity 实体, 用于存储实例的信息
type entity struct {
	alias       string          // 实体的名称,如果为空则使用其 reflect.Type 作为 key
	dependOn    map[string]bool // 实体依赖其他实体的列表
	t           reflect.Type    // 实体对应的类型
	v           reflect.Value   // 实体对应的值
	instance    interface{}     // 实体实例
	constructor interface{}     // 构造函数
	sync.RWMutex
}

// init 初始化实体
func (e *entity) init(val interface{}) *entity {
	t := reflect.TypeOf(val)
	v := reflect.ValueOf(val)
	e.t = t
	e.v = v
	e.dependOn = make(map[string]bool)
	e.instance = val
	st := utils.RemoveTypePtr(t)
	for i := 0; i < st.NumField(); i++ {
		f := st.Field(i)
		if _, ok := f.Tag.Lookup(injectTag); !ok {
			continue
		}
		fieldName := getFiledName(f)
		if _, exist := e.dependOn[fieldName]; !exist {
			e.dependOn[fieldName] = false
		}
	}
	return e
}

type Option func(*entity)

func WithAlias(alias string) Option {
	return func(o *entity) {
		o.alias = alias
	}
}

// canSet 检查所有字段是否全部设置完成
func (e *entity) canSet() bool {
	e.RLock()
	defer e.RUnlock()
	for _, v := range e.dependOn {
		if !v {
			return false
		}
	}
	return true
}

func (e *entity) setDeep(field reflect.StructField) {
	fn := getFiledName(field)
	e.Lock()
	if _, exist := e.dependOn[fn]; exist {
		e.dependOn[fn] = true
	}
	e.Unlock()
}

// getName 获取entity的名称
func (e *entity) getName() string {
	if len(e.alias) > 0 {
		return e.alias
	}
	return e.t.Name()
}

// getFiledName 获取字段的名称
func getFiledName(filed reflect.StructField) string {
	name, _ := filed.Tag.Lookup(injectTag)
	if len(name) > 0 {
		return name
	}
	return filed.Type.Name()
}
