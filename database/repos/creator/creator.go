package creator

import "context"

type Creator struct {
	Tx   context.Context // 事务上下文
	Data interface{}     // 保存的数据
}

func New() *Creator {
	return &Creator{}
}

func (c *Creator) SetTx(tx context.Context) *Creator {
	c.Tx = tx
	return c
}

func (c *Creator) SetData(data interface{}) *Creator {
	c.Data = data
	return c
}
