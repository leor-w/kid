package container

import (
	"errors"
	"reflect"
)

type bucket struct {
	named   map[string]*entity
	unnamed map[reflect.Type]*entity
}

var (
	ErrEntityExist = errors.New("entity existed")
)

func (b *bucket) add(e *entity) error {
	if len(e.name) > 0 {
		if _, exist := b.named[e.name]; exist {
			return ErrEntityExist
		}
		b.named[e.name] = e
		return nil
	}
	if _, exist := b.unnamed[e.t]; exist {
		return ErrEntityExist
	}
	b.unnamed[e.t] = e
	return nil
}

func (b *bucket) get(t reflect.Type, name string) *entity {
	if len(name) > 0 {
		return b.named[name]
	}
	return b.unnamed[t]
}

func newBucket() *bucket {
	return &bucket{
		named:   make(map[string]*entity),
		unnamed: make(map[reflect.Type]*entity),
	}
}
