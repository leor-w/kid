package container

import (
	"errors"
	"fmt"
	"github.com/leor-w/utils"
	"reflect"
	"sync/atomic"
	"unsafe"
)

const injectTag = "inject"

type Container interface {
	Provider
	Populater
	Invoker
}

type Provider interface {
	ProvideFunc(val interface{}, opts ...Option) error
	Provide(val interface{}, opts ...Option) error
}

type Populater interface {
	Populate() error
}

type Invoker interface {
	Invoke(fn interface{}) ([]reflect.Value, error)
}

var (
	ErrorDependNotCompleted = errors.New("依赖尚未完成")
)

type container struct {
	buckets      map[reflect.Type]*bucket
	populateChan chan *entity
	count        uint32
}

func (c *container) ProvideFunc(v interface{}, opts ...Option) error {
	t := reflect.TypeOf(v)
	if !utils.IsFunc(v) {
		return errors.New("container.ProvideFunc: not func or is nil")
	}
	for i := 0; i < t.NumOut(); i++ {
		e := newEntity(v)
		for _, o := range opts {
			o(e)
		}
	}
	return nil
}

func (c *container) Provide(i interface{}, opts ...Option) error {
	e := newEntity(i)
	for _, o := range opts {
		o(e)
	}
	if utils.IsNilPointer(i) {
		return errors.New("inject.Injector.Provide: 实例不是一个有效的指针值")
	}
	v := reflect.ValueOf(i)
	t := v.Type()
	b, exist := c.buckets[t]
	if !exist || b == nil {
		b = newBucket()
		c.buckets[t] = b
	}
	if err := b.add(e); err != nil {
		return err
	}
	atomic.AddUint32(&c.count, 1)
	return nil
}

func (c *container) Populate() error {
	c.populateChan = make(chan *entity, c.count)
	for _, b := range c.buckets {
		for _, e := range b.named {
			c.populateChan <- e
		}
		for _, e := range b.unnamed {
			c.populateChan <- e
		}
	}
	var err error
	for {
		if len(c.populateChan) <= 0 {
			break
		}
		select {
		case e := <-c.populateChan:
			// 如果依赖尚未完备 则将其置入队列末尾
			if err = c.populate(e); err != nil {
				if errors.Is(err, ErrorDependNotCompleted) {
					c.populateChan <- e
					continue
				}
				return fmt.Errorf("inject.Injector.Populate: %w", err)
			}
		default:
			break
		}
	}
	return nil
}

func (c *container) populate(e *entity) error {
	v := e.v
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil
	}
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		field := t.Field(i)
		tag, ok := field.Tag.Lookup(injectTag)
		if !ok {
			continue
		}
		if !f.CanSet() {
			continue
		}
		ft := f.Type()
		fe, err := c.Get(ft, tag)
		if err != nil {
			return err
		}
		if !fe.canSet() {
			return ErrorDependNotCompleted
		}
		if !fe.v.IsValid() {
			return fmt.Errorf("未找到 type: %v, name: %s", ft, e.name)
		}
		if !f.CanSet() {
			f = reflect.NewAt(f.Type(), unsafe.Pointer(f.UnsafeAddr())).Elem()
		}
		f.Set(fe.v)
		name := getFiledName(field)
		if _, exist := e.dependOn[name]; exist {
			e.dependOn[name] = true
		}
	}
	return nil
}

func (c *container) Invoke(fn interface{}) ([]reflect.Value, error) {
	ft := reflect.TypeOf(fn)
	if ft == nil || ft.Kind() != reflect.Func {
		return nil, errors.New("inject.Injector.Invoke: 必须是一个函数")
	}
	return nil, nil
}

func (c *container) Get(t reflect.Type, name string) (*entity, error) {
	b, exist := c.buckets[t]
	if !exist {
		return nil, fmt.Errorf("container.Get: not found bucket [%v]", t)
	}
	e := b.get(t, name)
	if e == nil {
		return nil, fmt.Errorf("container.Get: not found entity type [%v] name [%s]", t, name)
	}
	return e, nil
}

func New() Container {
	return &container{buckets: make(map[reflect.Type]*bucket)}
}
