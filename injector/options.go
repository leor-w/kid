package injector

import "reflect"

type Options struct {
	Alias string
	Scope string
}

type Option func(*Options)

func WithAlias(alias string) Option {
	return func(o *Options) {
		o.Alias = alias
	}
}

func WithScope(scope string) Option {
	return func(o *Options) {
		o.Scope = scope
	}
}

type InvokeOptions struct {
	InvokeInfo []*InvokeInfo
}

type InvokeInfo struct {
	t                      reflect.Type
	optional, alias, scope string
}

func NewInvokeInfo(v interface{}, scope, alias, optional string) *InvokeInfo {
	return &InvokeInfo{
		t:        reflect.TypeOf(v),
		optional: optional,
		alias:    alias,
		scope:    scope,
	}
}

type InvokeOption func(*InvokeOptions)

func WithInvokeInfo(info ...*InvokeInfo) InvokeOption {
	return func(o *InvokeOptions) {
		o.InvokeInfo = append(o.InvokeInfo, info...)
	}
}

func (i *InvokeOptions) GetInvokeInfo(t reflect.Type) *InvokeInfo {
	for _, info := range i.InvokeInfo {
		if info.t == t {
			return info
		}
	}
	return nil
}
