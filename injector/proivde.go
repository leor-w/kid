package injector

import "context"

type IProvider interface {
	Provide(ctx context.Context) any // 为依赖项提供值
}
