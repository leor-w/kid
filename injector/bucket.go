package injector

import (
	"reflect"
	"sync"
)

type bucket struct {
	// 指定名称的实例
	hasAlias map[string]*entity

	// 未指定名称的实例
	noAlias map[reflect.Type]*entity

	sync.RWMutex
}

func (b *bucket) set(e *entity) {
	if len(e.alias) > 0 {
		b.RLock()
		if _, exist := b.hasAlias[e.alias]; !exist {
			b.hasAlias[e.alias] = e
		}
		b.RUnlock()
		return
	}
	if _, exist := b.noAlias[e.t]; !exist {
		b.RLock()
		b.noAlias[e.t] = e
		b.RUnlock()
	}
}

func (b *bucket) get(t reflect.Type, alias string) *entity {
	b.RLock()
	defer b.RUnlock()
	if len(alias) > 0 {
		return b.hasAlias[alias]
	}
	return b.noAlias[t]
}

func newBucket() *bucket {
	return &bucket{
		hasAlias: make(map[string]*entity),
		noAlias:  make(map[reflect.Type]*entity),
	}
}
