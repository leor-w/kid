package plugin

import "context"

type Plugin interface {
	Provide(context.Context) interface{}
}

type NameKey struct{}
