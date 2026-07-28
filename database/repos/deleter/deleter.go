package deleter

import (
	"context"

	"github.com/lsjhtang/kid/database/repos/where"
)

type Deleter struct {
	// 事务
	Tx context.Context

	// 模型
	Model interface{}

	// 是否使用软删除, 默认使用软删除, true: 使用硬删除
	Unscoped bool

	// 是否开启调试模式
	Debug bool

	// 关联
	Associations []string

	// 是否关联所有
	IsAssociationsAll bool

	// 条件
	Wheres *where.Wheres
}

func New() *Deleter {
	return &Deleter{}
}

func (d *Deleter) SetTx(tx context.Context) *Deleter {
	d.Tx = tx
	return d
}

func (d *Deleter) SetModel(model interface{}) *Deleter {
	d.Model = model
	return d
}

func (d *Deleter) SetWheres(wheres *where.Wheres) *Deleter {
	d.Wheres = wheres
	return d
}
