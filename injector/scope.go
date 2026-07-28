package injector

import (
	"fmt"
	"reflect"
	"strings"
	"sync/atomic"
	"unsafe"

	"github.com/lsjhtang/kid/utils"
)

type Scope struct {
	// name 当前 Scope 的名称
	name string
	// 存储所有实例的桶
	buckets map[reflect.Type]*bucket

	// 实例数量
	count uint32

	// parentScope 当前 Scope 的父 Scope
	parentScope *Scope

	// childScopes 当前 Scope 的子 Scope
	childScopes []*Scope
}

func (scope *Scope) Scope(name string) *Scope {
	var child = &Scope{
		name:        name,
		buckets:     make(map[reflect.Type]*bucket),
		childScopes: make([]*Scope, 0),
	}
	if scope != nil {
		scope.childScopes = append(scope.childScopes, child)
	}
	child.parentScope = scope
	return child
}

// rootScope 返回当前 Scope 的根 Scope
func (scope *Scope) rootScope() *Scope {
	if scope.parentScope == nil {
		return scope
	}
	return scope.parentScope.rootScope()
}

func (scope *Scope) getScope(name string) *Scope {
	if scope.name == name {
		return scope
	}
	for _, child := range scope.childScopes {
		if child.name == name {
			return child
		}
		s := child.getScope(name)
		if s != nil {
			return s
		}
	}
	return nil
}

func (scope *Scope) Provide(provider IProvider, option ...Option) error {
	var options = &Options{}
	for _, o := range option {
		o(options)
	}
	return scope.provide(provider, options)
}

func (scope *Scope) provide(provider IProvider, options *Options) error {
	e := newEntity(provider, options)
	if utils.IsNilPointer(e.instance) {
		return fmt.Errorf("container.Provide: 实例为无效的空指针")
	}
	scope.setValue(e.t, e)
	return nil
}

func (scope *Scope) Invoke(fn interface{}) ([]reflect.Value, error) {
	v := reflect.ValueOf(fn)
	if v.Kind() != reflect.Func {
		return nil, fmt.Errorf("container.Invoke: 必须是一个函数")
	}
	t := v.Type()
	in := make([]reflect.Value, t.NumIn())
	for i := 0; i < t.NumIn(); i++ {
		it := t.In(i)
		e, err := scope.get(it, nil)
		if err != nil {
			return nil, err
		}
		if !e.isComplete() {
			return nil, fmt.Errorf(InjectionUnfinishedError, it)
		}
		if !e.v.IsValid() {
			return nil, fmt.Errorf(InvalidInjectionFiledError, it, e.alias)
		}
		in[i] = e.v
	}
	return v.Call(in), nil
}

func (scope *Scope) getPopulateChan() chan *entity {
	popChan := make(chan *entity, scope.count)
	for _, bucket := range scope.buckets {
		for _, e := range bucket.hasAlias {
			popChan <- e
		}
		for _, e := range bucket.noAlias {
			popChan <- e
		}
	}
	return popChan
}

func (scope *Scope) Populate() error {
	return scope.populate()
}

// populate 依赖注入
func (scope *Scope) populate() error {
	popChan := scope.getPopulateChan()
	for {
		if len(popChan) == 0 {
			break
		}
		select {
		case e := <-popChan:
			if err := scope.popEntity(e); err != nil {
				return err
			}
			if !e.isComplete() {
				e.printEntityDependy()
				popChan <- e
				continue
			}
		default:
			break
		}
	}

	if len(scope.childScopes) > 0 {
		for _, child := range scope.childScopes {
			if err := child.populate(); err != nil {
				return fmt.Errorf("[%s] %w", child.name, err)
			}
		}
	}
	return nil
}

func (scope *Scope) popEntity(e *entity) error {
	v := e.v
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return fmt.Errorf(NonPointerError, v)
	}
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		field := t.Field(i)
		//if e.isDependency(field) {
		//	continue
		//}
		// 检查是否有注入标签
		tag, ok := field.Tag.Lookup(injectTag)
		if !ok {
			continue
		}
		// 通过注入标签获取对应类型的注入实例
		ft := f.Type()
		tm := parseTag(tag)
		fe, err := scope.get(ft, tm)
		if err != nil {
			return fmt.Errorf(NotFoundEntityError, v.Type(), ft, scope.name, tm.getAlias())
		}
		if !fe.v.IsValid() {
			return fmt.Errorf(InvalidInjectionFiledError, ft, e.alias)
		}
		if !f.CanSet() {
			f = reflect.NewAt(f.Type(), unsafe.Pointer(f.UnsafeAddr())).Elem()
		}
		f.Set(fe.v)
		e.setDependency(field)
	}
	return nil
}

func (scope *Scope) getRecursive(t reflect.Type, tm tagMapper) (*entity, error) {
	e, err := scope.get(t, tm)
	if err != nil || !e.isComplete() || !e.v.IsValid() {
		if len(scope.childScopes) <= 0 {
			return nil, err
		}
		for _, child := range scope.childScopes {
			e, err = child.getRecursive(t, tm)
			if err != nil || e == nil || !e.isComplete() || !e.v.IsValid() {
				continue
			}
			return e, nil
		}
	}
	return e, nil
}

func (scope *Scope) get(t reflect.Type, tag tagMapper) (*entity, error) {
	tagScope := scope
	if tag.hasScope() {
		tagScope = scope.rootScope().getScope(tag.getScope())
	}
	bt, b, err := tagScope.tryGetBucket(t, tag)
	if err != nil || b == nil {
		// 找不到到, 检查是否标记为可选 如果是则返回 nil
		if strings.Contains(tag.getOptional(), TagOptionNotRequired) {
			return nil, nil
		}
		return nil, err
	}
	if bt != nil {
		t = bt
	}

	// 获取对应的 entity
	e := b.get(t, tag.getAlias())
	if e == nil {
		return nil, fmt.Errorf("未找到对应的 entity 类型为: [%v] 别名为: [%s]", t, tag.getAlias())
	}
	return e, nil
}

// tryGetBucket 尝试获取对应的 bucket
func (scope *Scope) tryGetBucket(t reflect.Type, tag tagMapper) (reflect.Type, *bucket, error) {
	// 直接找对应的 bucket
	b, exist := scope.buckets[t]
	if exist {
		return t, b, nil
	}
	bt, b, err := scope.findAssignableBucket(t, tag)
	if err != nil {
		return scope.findParentScopeBucket(t, tag)
	}
	return bt, b, nil
}

// findAssignableBucket 查找可赋值或为接口实现的 bucket
func (scope *Scope) findAssignableBucket(t reflect.Type, tag tagMapper) (reflect.Type, *bucket, error) {
	for k, v := range scope.buckets {
		if k.AssignableTo(t) || t.AssignableTo(k) {
			if tag != nil && len(tag.getAlias()) > 0 {
				ve := v.get(t, tag.getAlias())
				if ve != nil {
					return k, v, nil
				}
			} else {
				return k, v, nil
			}
		}
	}
	return nil, nil, fmt.Errorf(NotFoundBucketError, t)
}

// findParentScopeBucket 从父容器中查找可赋值或为接口实现的 bucket
func (scope *Scope) findParentScopeBucket(t reflect.Type, tag tagMapper) (reflect.Type, *bucket, error) {
	if scope.parentScope == nil {
		return nil, nil, fmt.Errorf(NotFoundBucketError, t)
	}
	pt, pb, _ := scope.parentScope.tryGetBucket(t, tag)
	if pb != nil {
		if pt != nil {
			t = pt
		}
		return pt, pb, nil
	}
	return nil, nil, fmt.Errorf(NotFoundBucketError, t)
}

// setValue 设置实例
func (scope *Scope) setValue(k reflect.Type, e *entity) {
	bucket, exist := scope.buckets[k]
	if !exist {
		bucket = newBucket()
		scope.buckets[k] = bucket
	}
	bucket.set(e)
	atomic.AddUint32(&scope.count, 1)
}
