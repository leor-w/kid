package finder

import (
	"context"
	"github.com/lsjhtang/kid/database/repos/where"
	"gorm.io/gorm"
)

type Finder struct {
	// 数据库事务上下文
	Tx context.Context
	// 查询的模型
	Model interface{}

	// 查询条件
	Wheres *where.Wheres

	// 查询结果保存的对象
	Recipient interface{}

	// 排序条件
	OrderBy string

	// 是否过滤未找到数据的错误
	IgnoreNotFound bool

	// 预加载关联数据的表名
	Preloads []string

	// 查询作用域
	Scopes []func(db *gorm.DB) *gorm.DB

	// 是否查询软删除的数据, 默认不查询软删除的数据, true: 查询软删除的数据
	Unscoped bool

	// 是否开启调试模式
	Debug bool

	// 分页页码, 单页大小
	Num, Size int

	// 分页总数
	Total *int64
}

// Sum 查询某个字段的和
type Sum struct {
	Model  interface{}
	Wheres *where.Wheres
	Col    string
	Val    interface{}
	Debug  bool
}

func New() *Finder {
	return &Finder{}
}

func (f *Finder) SetModel(model interface{}) *Finder {
	f.Model = model
	return f
}

func (f *Finder) SetWheres(wheres *where.Wheres) *Finder {
	f.Wheres = wheres
	return f
}

func (f *Finder) SetRecipient(recipient interface{}) *Finder {
	f.Recipient = recipient
	return f
}

func (f *Finder) SetOrderBy(orderBy string) *Finder {
	f.OrderBy = orderBy
	return f
}

func (f *Finder) SetIgnoreNotFound(ignoreNotFound bool) *Finder {
	f.IgnoreNotFound = ignoreNotFound
	return f
}

func (f *Finder) SetPreloads(preloads ...string) *Finder {
	f.Preloads = preloads
	return f
}

func (f *Finder) SetScope(scopes ...func(db *gorm.DB) *gorm.DB) *Finder {
	f.Scopes = scopes
	return f
}

func (f *Finder) SetPaginate(num, size int, total *int64) *Finder {
	f.Num = num
	f.Size = size
	f.Total = total
	return f
}
