package repos

import (
	"context"

	"github.com/lsjhtang/kid/database/repos/deleter"

	"github.com/lsjhtang/kid/database/repos/updater"

	"github.com/lsjhtang/kid/database/repos/creator"

	"github.com/lsjhtang/kid/database/repos/finder"
)

type (
	// IBasicRepository MySQL基础查询与事务
	IBasicRepository interface {
		IRepository
		ITxRepository
	}
	// IRepository MySQL基础查询
	IRepository interface {
		Exist(*finder.Finder) bool
		GetOne(*finder.Finder) error
		GetById(id int64, model interface{}) error
		GetByKV(kv map[string]interface{}, model interface{}) error
		Find(*finder.Finder) error
		Create(*creator.Creator) error
		Update(*updater.Updater) error
		UpdateByKV(data interface{}, kv map[string]interface{}) error
		Delete(*deleter.Deleter) error
		Save(*creator.Creator) error
		Exec(context.Context, string, ...interface{}) error
		Raw(context.Context, string, ...interface{}) error
		GetUniqueID(finder *finder.Finder, min, max, ignoreStart, ignoreEnd int64) int64
		Count(*finder.Finder) error
		Sum(*finder.Sum) error
		GetTableName(data any) string
		CreateTable(tableName string, data any) error
		TableHasExist(tableName string) bool
	}
	// ITxRepository MySQL事务
	ITxRepository interface {
		Transaction(func(context.Context) error) error                             // 开启新的事物执行
		GetDb(context.Context) interface{}                                         // 获取数据库连接
		WhetherTx(context.Context) bool                                            // 检查是否在事务中
		ExecWithTx(hasTx context.Context, fn func(tx context.Context) error) error // 如果hasTx已经在事务中,则不开启新的事务, 否则开启新的事务执行 fn
	}
	// IRedisRepository Redis基础查询
	IRedisRepository interface {
		Set(key string, value interface{}, expire int64) error
		Get(key string) (string, error)
		Del(keys ...string) error
		Expire(key string, expire int64) error
		Exists(key string) (bool, error)
	}
)

type TxKey struct{}
