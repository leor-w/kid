package container

import (
	"github.com/leor-w/utils"
	"reflect"
)

type entity struct {
	name        string          // 实体的名称,如果为空则使用其 reflect.Type 作为 key
	dependOn    map[string]bool // 实体依赖其他实体的列表
	t           reflect.Type    // 实体对应的类型
	v           reflect.Value   // 实体对应的值
	instance    interface{}     // 实体实例
	constructor interface{}     // 构造函数
}

func newEntity(i interface{}) *entity {
	t := reflect.TypeOf(i)
	v := reflect.ValueOf(i)
	e := entity{
		t:        t,
		v:        v,
		dependOn: make(map[string]bool),
	}
	//if tools.IsFunc(i) {
	//	e.constructor = i
	//	for i := 0; i < t.NumIn(); i++ {
	//		if _, exist := e.dependOn[t.In(i)]; !exist {
	//			e.dependOn[t.In(i)] = false
	//		}
	//	}
	//	return &e
	//}
	e.instance = i
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
	return &e
}

type Option func(*entity)

func WithName(name string) Option {
	return func(o *entity) {
		o.name = name
	}
}

func (e *entity) canSet() bool {
	for _, v := range e.dependOn {
		if !v {
			return false
		}
	}
	return true
}

func (e *entity) getName() string {
	if len(e.name) > 0 {
		return e.name
	}
	return e.t.Name()
}

func getFiledName(filed reflect.StructField) string {
	name, _ := filed.Tag.Lookup(injectTag)
	if len(name) > 0 {
		return name
	}
	return filed.Type.Name()
}
