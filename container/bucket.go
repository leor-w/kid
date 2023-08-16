package container

import (
	"reflect"
	"sync"
)

// bucket 桶, 用于存储所有实例
type bucket struct {
	sync.RWMutex
	named   map[string]*entity       // 命名实例
	unnamed map[reflect.Type]*entity // 未命名实例
}

func (b *bucket) add(e *entity) {
	if len(e.alias) > 0 {
		b.RLock()
		if _, exist := b.named[e.alias]; !exist {
			b.named[e.alias] = e
		}
		b.RUnlock()
		return
	}
	if _, exist := b.unnamed[e.t]; !exist {
		b.RLock()
		b.unnamed[e.t] = e
		b.RUnlock()
	}
}

func (b *bucket) getEntity(t reflect.Type, name string) *entity {
	if len(name) > 0 {
		return b.named[name]
	}
	b.RLock()
	defer b.RUnlock()
	return b.unnamed[t]
}

func newBucket() *bucket {
	return &bucket{
		named:   make(map[string]*entity),
		unnamed: make(map[reflect.Type]*entity),
	}
}
