package updater

import (
	"context"

	"github.com/lsjhtang/kid/database/repos/where"
)

type Updater struct {
	Tx       context.Context        // 事务上下文
	Model    interface{}            // 修改对应的模型
	Wheres   *where.Wheres          // 修改的过滤条件
	Update   interface{}            // 修改的结构体值, 修改整体值时使用
	Fields   map[string]interface{} // 修改的map值, 用于修改单个或多个字段时使用
	SaveNil  bool                   // 是否保存空值, 默认不会更新类似 0,false,"" 等值
	SaveFull bool                   // 是否同时修改关联数据
	Debug    bool                   // 是否开启调试模式
	Preloads []string               // 预加载关联数据
	Omits    []string               // 修改时指定不更新的字段
	Selects  []string               // 修改时指定更新的字段
}

func New() *Updater {
	return &Updater{}
}

func (u *Updater) SetTx(tx context.Context) *Updater {
	u.Tx = tx
	return u
}

func (u *Updater) SetModel(model interface{}) *Updater {
	u.Model = model
	return u
}

func (u *Updater) SetWheres(wheres *where.Wheres) *Updater {
	u.Wheres = wheres
	return u
}

func (u *Updater) SetUpdate(update interface{}) *Updater {
	u.Update = update
	return u
}

func (u *Updater) SetFields(fields map[string]interface{}) *Updater {
	u.Fields = fields
	return u
}

func (u *Updater) SetSaveNil(saveNil bool) *Updater {
	u.SaveNil = saveNil
	return u
}

func (u *Updater) SetSaveFull(saveFull bool) *Updater {
	u.SaveFull = saveFull
	return u
}

func (u *Updater) SetPreloads(preloads ...string) *Updater {
	u.Preloads = preloads
	return u
}

func (u *Updater) SetOmits(omits []string) *Updater {
	u.Omits = omits
	return u
}
