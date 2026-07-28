package injector

import (
	"context"
	"reflect"
	"sync"

	"github.com/lsjhtang/kid/utils"
)

// entity 实体, 用于存储实例的信息
type entity struct {
	alias       string          // 实体的名称,如果为空则使用其 reflect.Type 作为 key
	scope       string          // 实体的作用域
	dependOn    map[string]bool // 实体依赖其他实体的列表
	t           reflect.Type    // 实体对应的类型
	v           reflect.Value   // 实体对应的值
	instance    interface{}     // 实体实例
	constructor interface{}     // 构造函数
	sync.RWMutex
}

// init 初始化实体
func newEntity(val IProvider, opts *Options) *entity {
	e := new(entity)
	e.alias = opts.Alias
	e.scope = opts.Scope
	ctx := context.WithValue(context.Background(), NameKey{}, e.alias)
	provide := val.Provide(ctx)
	return e.init(provide)
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

// isComplete 检查所有字段是否全部设置完成
func (e *entity) isComplete() bool {
	e.RLock()
	defer e.RUnlock()
	for _, v := range e.dependOn {
		if !v {
			return false
		}
	}
	return true
}

func (e *entity) printEntityDependy() {
	for k, v := range e.dependOn {
		if v {
			println(k, "已完成依赖")
		} else {
			println(k, "尚未完成依赖")
		}
	}
}

// setDependency 设置依赖
func (e *entity) setDependency(field reflect.StructField) {
	fn := getFiledName(field)
	e.Lock()
	if _, exist := e.dependOn[fn]; exist {
		e.dependOn[fn] = true
	}
	e.Unlock()
}

func (e *entity) isDependency(field reflect.StructField) bool {
	fn := getFiledName(field)
	e.RLock()
	defer e.RUnlock()
	if depend, exist := e.dependOn[fn]; exist {
		if depend {
			return true
		}
		return false
	}
	return false
}

// getFiledName 获取字段的名称
func getFiledName(filed reflect.StructField) string {
	name, _ := filed.Tag.Lookup(injectTag)
	if len(name) > 0 {
		return name
	}
	name = filed.Type.Name()
	if len(name) > 0 {
		return name
	}
	return filed.Name
}
