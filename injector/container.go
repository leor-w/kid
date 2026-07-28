package injector

import (
	"fmt"
	"reflect"

	"github.com/lsjhtang/kid/utils"
)

var container *Container

func init() {
	container = New()
}

func Provide(val IProvider, opts ...Option) error {
	return container.Provide(val, opts...)
}

func Invoke(fn interface{}, opts ...InvokeOption) ([]reflect.Value, error) {
	return container.Invoke(fn, opts...)
}

func Populate(opts ...Option) error {
	return container.Populate(opts...)
}

func SubScope(name string) *Scope {
	return container.Scope(name)
}

type Container struct {
	scope *Scope // 存储所有的 Scope
}

func (container *Container) Provide(val IProvider, opts ...Option) error {
	var options = &Options{}
	for _, o := range opts {
		o(options)
	}
	scope := container.Scope(options.Scope)
	return scope.provide(val, options)
}

func (container *Container) Invoke(fn interface{}, opts ...InvokeOption) ([]reflect.Value, error) {
	var options = &InvokeOptions{}
	for _, o := range opts {
		o(options)
	}
	if !utils.IsFunc(fn) {
		return nil, fmt.Errorf("container.Invoke: fn 必须是一个函数")
	}
	f := reflect.ValueOf(fn)
	t := f.Type()
	inArgs := make([]reflect.Value, t.NumIn())
	for i := 0; i < t.NumIn(); i++ {
		inArg := t.In(i)
		info := options.GetInvokeInfo(inArg)
		var (
			scope *Scope
			tm    tagMapper
		)
		if info != nil {
			scope = container.Scope(info.scope)
			tm = newTagMapper(info.alias, info.scope, info.optional)
		} else {
			scope = container.scope
		}
		e, err := scope.getRecursive(inArg, tm)
		if err != nil {
			return nil, err
		}
		if e == nil {
			return nil, fmt.Errorf("container.Invoke: 未找到 %s 的实例", inArg.String())
		}
		if !e.isComplete() {
			return nil, fmt.Errorf(InjectionUnfinishedError, inArg)
		}
		if !e.v.IsValid() {
			return nil, fmt.Errorf(InvalidInjectionFiledError, inArg, e.alias)
		}
		inArgs[i] = e.v
	}
	return f.Call(inArgs), nil
}

func (container *Container) Populate(opts ...Option) error {
	var options = &Options{}
	for _, o := range opts {
		o(options)
	}
	scope := container.Scope(options.Scope)
	if err := scope.populate(); err != nil {
		return err
	}
	return nil
}

func (container *Container) Scope(name string) *Scope {
	if name == "" {
		return container.scope
	}
	scope := container.scope.getScope(name)
	if scope != nil {
		return scope
	}
	return container.scope.Scope(name)
}

func New() *Container {
	return &Container{
		scope: &Scope{
			buckets:     make(map[reflect.Type]*bucket),
			childScopes: make([]*Scope, 0),
		},
	}
}
