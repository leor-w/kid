package container

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/leor-w/kid/plugin"
	"github.com/leor-w/kid/utils"
)

var (
	ErrorDependNotCompleted = errors.New("依赖尚未完成")
)

// Container 容器实现
type Container struct {
	// 父容器
	parent *Container

	buckets map[reflect.Type]*bucket // 存储所有实例的桶, key为实例类型, value为桶, 桶内存储了所有实例, 以及实例的依赖关系
	count   uint32                   // 实例数量

	sync.Mutex
}

func (c *Container) NewChild() *Container {
	child := New()
	child.parent = c
	return child
}

func (c *Container) ProvideFunc(v interface{}, opts ...Option) error {
	return nil
}

func (c *Container) Provide(plug plugin.Plugin, opts ...Option) error {
	e := c.initEntity(plug, opts...)
	if utils.IsNilPointer(e.instance) {
		return errors.New("container.Provide: 实例不是一个有效的指针值")
	}
	v := reflect.ValueOf(e.instance)
	t := v.Type()
	b, exist := c.buckets[t]
	if !exist || b == nil {
		b = newBucket()
		c.buckets[t] = b
	}
	b.add(e)
	atomic.AddUint32(&c.count, 1)
	return nil
}

func (c *Container) initEntity(plug plugin.Plugin, opts ...Option) *entity {
	e := new(entity)
	for _, o := range opts {
		o(e)
	}
	ctx := context.WithValue(context.Background(), plugin.NameKey{}, e.alias)
	provide := plug.Provide(ctx)
	return e.init(provide)
}

func (c *Container) initPopulateCh(ch chan<- *entity) {
	for _, b := range c.buckets {
		for _, e := range b.named {
			ch <- e
		}
		for _, e := range b.unnamed {
			ch <- e
		}
	}
	if c.parent != nil {
		c.parent.initPopulateCh(ch)
	}
}

func (c *Container) Populate() error {
	if c.parent != nil {
		if err := c.parent.Populate(); err != nil {
			return err
		}
	}
	populateCh := make(chan *entity, c.count)
	for _, b := range c.buckets {
		for _, e := range b.named {
			populateCh <- e
		}
		for _, e := range b.unnamed {
			populateCh <- e
		}
	}
	for {
		if len(populateCh) <= 0 {
			break
		}
		select {
		case e := <-populateCh:
			// 如果依赖尚未完备 则将其置入队列末尾
			complete, err := c.populate(e)
			if err != nil {
				if errors.Is(err, ErrorDependNotCompleted) {
					populateCh <- e
					continue
				}
				return fmt.Errorf("container.Populate: %w", err)
			}
			if !complete {
				populateCh <- e
				continue
			}
		default:
			break
		}
	}
	return nil
}

func (c *Container) PopulateSingle(val plugin.Plugin, opts ...Option) (bool, error) {
	if utils.IsNilPointer(val) {
		return false, errors.New("inject.Injector.PopulateSingle: 实例不是一个有效的指针值")
	}
	e := c.initEntity(val, opts...)

	return c.populate(e)
}

func (c *Container) populate(e *entity) (bool, error) {
	comp := true
	v := e.v
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return false, fmt.Errorf("Container.populate: 注入类型 [%v] 错误, 必须为一个指针", v)
	}
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		field := t.Field(i)
		tag, ok := field.Tag.Lookup(injectTag)
		if !ok {
			continue
		}
		it := parseTag(tag)
		ft := f.Type()
		fe, err := c.get(ft, it)
		if err != nil {
			return false, fmt.Errorf("container.populate: [%v] 未找到对应的实例: [%v]", v.Type(), ft)
		}
		if fe == nil && err == nil {
			e.setDeep(field)
			return comp, nil
		}
		if fe == nil {
			return false, fmt.Errorf("Container.populate: [%v] 未找到对应的实例: [%v]", v, ft)
		}
		if !fe.canSet() {
			comp = false
		}
		if !fe.v.IsValid() {
			return false, fmt.Errorf("container.populate: 类型 [%v], 别名: [%s] 值无效", ft, e.alias)
		}
		if !f.CanSet() {
			f = reflect.NewAt(f.Type(), unsafe.Pointer(f.UnsafeAddr())).Elem()
		}
		f.Set(fe.v)
		e.setDeep(field)
	}
	return comp, nil
}

func (c *Container) Invoke(fn interface{}) ([]reflect.Value, error) {
	ft := reflect.TypeOf(fn)
	if ft == nil || ft.Kind() != reflect.Func {
		return nil, errors.New("container.Invoke: 必须是一个函数")
	}
	return nil, nil
}

func (c *Container) Get(plugin plugin.Plugin, names ...string) (interface{}, error) {
	var name string
	if len(names) > 0 {
		name = names[0]
	}
	entity, err := c.get(reflect.TypeOf(plugin), InjectTag{name: name})
	if err != nil {
		return nil, err
	}
	return entity.instance, nil
}

func (c *Container) get(t reflect.Type, it InjectTag) (*entity, error) {
	bt, b, err := c.getBucket(t)
	if err != nil || b == nil {
		// 找不到到, 检查是否标记为可选 如果是则返回 nil
		if strings.Contains(it.getOpts(), TagOptionNotRequired) {
			return nil, nil
		}
		return nil, err
	}
	if bt != nil {
		t = bt
	}

	// 获取对应的 entity
	e := b.getEntity(t, it.getAlias())
	if e == nil {
		return nil, fmt.Errorf("container.getEntity: 未找到对应的 entity 类型为: [%v] 别名为: [%s]", t, it.getAlias())
	}
	return e, nil
}

func (c *Container) getBucket(t reflect.Type) (reflect.Type, *bucket, error) {
	// 直接找对应的 bucket
	b, exist := c.buckets[t]
	if !exist {
		// 如果无法直接找到对应的 bucket 则遍历所有 bucket 寻找可赋值的类型 例如: 接口 -> 实现
		for k, v := range c.buckets {
			if k.AssignableTo(t) || t.AssignableTo(k) {
				//t = k
				//b = v
				return k, v, nil
			}
		}

		// 如果没有找到 bucket 则检查是否有父容器
		if c.parent == nil {
			return nil, nil, fmt.Errorf("container.getBucket: 未找到对应的 bucket: [%v]", t)
		}
		// 从父容器中查找
		pt, pb, _ := c.parent.getBucket(t)
		if pb != nil {
			if pt != nil {
				t = pt
			}
			return t, pb, nil
		}
		// 如果父容器中也没有找到则返回错误
		return nil, nil, fmt.Errorf("container.getBucket: 未找到对应的 bucket: [%v]", t)
	}
	return nil, b, nil
}

func New() *Container {
	return &Container{
		buckets: make(map[reflect.Type]*bucket),
	}
}

// NewWithParent 创建一个新的容器, 并指定父容器
func NewWithParent(parent *Container) *Container {
	return &Container{
		parent:  parent,
		buckets: make(map[reflect.Type]*bucket),
	}
}
