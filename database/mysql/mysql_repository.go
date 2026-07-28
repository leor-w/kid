package mysql

import (
	"context"
	"errors"
	"reflect"

	"gorm.io/gorm/clause"

	"github.com/lsjhtang/kid/database/repos"
	"github.com/lsjhtang/kid/database/repos/creator"
	"github.com/lsjhtang/kid/database/repos/deleter"
	"github.com/lsjhtang/kid/database/repos/finder"
	"github.com/lsjhtang/kid/database/repos/updater"
	"github.com/lsjhtang/kid/logger"
	"github.com/lsjhtang/kid/utils"
	"gorm.io/gorm"
)

var (
	MissingUpdates = errors.New("更新内容为空")
)

type Repository struct {
	DB *MySQL `inject:""`
}

func (repo *Repository) Provide(context.Context) interface{} {
	return repo
}

func (repo *Repository) handleErr(ignoreErr bool, err error) error {
	if ignoreErr && errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Debug("record not found")
		return nil
	} else {
		return err
	}
}

func (repo *Repository) Transaction(txFunc func(context.Context) error) error {
	tx := repo.DB.Begin()
	ctx := context.WithValue(context.Background(), repos.TxKey{}, tx)
	err := txFunc(ctx)
	if err != nil {
		if err := tx.Rollback().Error; err != nil {
			return err
		}
		return err
	}
	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
}

func (repo *Repository) GetDb(ctx context.Context) interface{} {
	return repo.getTx(ctx)
}

func (repo *Repository) WhetherTx(tx context.Context) bool {
	return repo.getTx(tx) != nil
}

func (repo *Repository) ExecWithTx(tx context.Context, fn func(tx context.Context) error) error {
	if repo.WhetherTx(tx) {
		return fn(tx)
	}
	return repo.Transaction(func(tx context.Context) error {
		return fn(tx)
	})
}

func (repo *Repository) getTx(tx context.Context) *gorm.DB {
	if tx == nil {
		return repo.DB.DB
	}
	txDb, ok := tx.Value(repos.TxKey{}).(*gorm.DB)
	if !ok {
		return repo.DB.DB
	}
	return txDb
}

func (repo *Repository) Exist(finder *finder.Finder) bool {
	var db = repo.getTx(finder.Tx)
	if finder.Debug {
		db.Debug()
	}
	var count int64
	if err := db.
		Model(finder.Model).
		Scopes(Wheres(finder.Wheres.Wheres...)).
		Count(&count).
		Error; err != nil {
		logger.Errorf("mysql_repository.go Exist error: %v", err.Error())
		return false
	}
	if count <= 0 {
		return false
	}
	return true
}

func (repo *Repository) GetOne(finder *finder.Finder) error {
	var db = repo.getTx(finder.Tx).Scopes(Wheres(finder.Wheres.Wheres...))
	if finder.Debug {
		db.Debug()
	}
	if len(finder.Preloads) > 0 {
		for _, preload := range finder.Preloads {
			db.Preload(preload)
		}
	}
	if len(finder.Scopes) > 0 {
		db.Scopes(finder.Scopes...)
	}
	if finder.Unscoped {
		db.Unscoped()
	}
	if err := db.First(finder.Recipient).Error; err != nil {
		return repo.handleErr(finder.IgnoreNotFound, err)
	}
	return nil
}

func (repo *Repository) GetById(id int64, model interface{}) error {
	return repo.DB.First(model, id).Error
}

func (repo *Repository) GetByKV(kv map[string]interface{}, model interface{}) error {
	return repo.DB.Where(kv).First(model).Error
}

func (repo *Repository) Find(finder *finder.Finder) error {
	var db = repo.getTx(finder.Tx)
	if finder.Wheres != nil && len(finder.Wheres.Wheres) > 0 {
		db = db.Scopes(Wheres(finder.Wheres.Wheres...))
	}
	if finder.Debug {
		db = db.Debug()
	}
	if finder.Model != nil {
		db = db.Model(finder.Model)
	}
	if finder.OrderBy != "" {
		db = db.Order(finder.OrderBy)
	}
	if len(finder.Preloads) > 0 {
		for _, preload := range finder.Preloads {
			db = db.Preload(preload)
		}
	}
	if len(finder.Scopes) > 0 {
		db = db.Scopes(finder.Scopes...)
	}
	if finder.Unscoped {
		db = db.Unscoped()
	}
	if finder.Total != nil {
		if err := db.Count(finder.Total).Error; err != nil {
			return repo.handleErr(finder.IgnoreNotFound, err)
		}
	}
	if finder.Size > 0 {
		db = db.Scopes(Paginate(finder.Num, finder.Size))
	}
	if err := db.Find(finder.Recipient).Error; err != nil {
		return repo.handleErr(finder.IgnoreNotFound, err)
	}
	return nil
}

func (repo *Repository) Create(creator *creator.Creator) error {
	var db = repo.getTx(creator.Tx)
	if creator.Debug {
		db.Debug()
	}
	return db.
		Create(creator.Data).
		Error
}

func (repo *Repository) Save(creator *creator.Creator) error {
	var db = repo.getTx(creator.Tx)
	if creator.Debug {
		db.Debug()
	}
	return db.Save(creator.Data).Error
}

func (repo *Repository) Update(updater *updater.Updater) error {
	var db = repo.getTx(updater.Tx)
	if updater.Wheres != nil && len(updater.Wheres.Wheres) > 0 {
		db = db.Scopes(Wheres(updater.Wheres.Wheres...))
	}
	if updater.Debug {
		db.Debug()
	}
	if len(updater.Omits) > 0 {
		db.Omit(updater.Omits...)
	}
	if updater.SaveNil && updater.Update != nil && len(updater.Selects) <= 0 {
		db = db.Select("*")
	}
	if len(updater.Selects) > 0 {
		db = db.Select(updater.Selects)
	}
	if updater.Update != nil {
		return db.
			Updates(updater.Update).
			Error
	}
	if updater.Fields != nil && updater.Model != nil {
		return db.Model(updater.Model).
			UpdateColumns(updater.Fields).
			Error
	}
	return MissingUpdates
}

func (repo *Repository) UpdateByKV(data interface{}, kv map[string]interface{}) error {
	return repo.DB.Model(data).Where(kv).Updates(data).Error
}

func (repo *Repository) Delete(deleter *deleter.Deleter) error {
	var db = repo.getTx(deleter.Tx)
	if deleter.Debug {
		db.Debug()
	}
	if deleter.IsAssociationsAll {
		db.Select(clause.Associations)
	}
	if len(deleter.Associations) > 0 {
		db.Select(deleter.Associations)
	}
	return db.
		Model(deleter.Model).
		Scopes(Wheres(deleter.Wheres.Wheres...)).
		Delete(nil).
		Error
}

func (repo *Repository) Exec(tx context.Context, sql string, args ...interface{}) error {
	return repo.getTx(tx).Exec(sql, args...).Error
}

func (repo *Repository) Raw(tx context.Context, sql string, args ...interface{}) error {
	return repo.getTx(tx).Raw(sql, args...).Error
}

func (repo *Repository) GetUniqueID(finder *finder.Finder, min, max, ignoreStart, ignoreEnd int64) int64 {
	var uniqueId int64
	for {
		id := utils.UniqueId(min, max)
		if ignoreStart > 0 && ignoreEnd > 0 && ignoreStart <= uniqueId && uniqueId <= ignoreEnd {
			continue
		}
		finder.Wheres.Wheres[0].V1 = id
		if repo.Exist(finder) {
			continue
		}
		uniqueId = id
		break
	}
	return uniqueId
}

func (repo *Repository) Count(finder *finder.Finder) error {
	var db = repo.DB
	if finder.Debug {
		db.Debug()
	}
	if err := db.
		Model(finder.Model).
		Scopes(Wheres(finder.Wheres.Wheres...)).
		Count(finder.Total).
		Error; err != nil {
		return repo.handleErr(finder.IgnoreNotFound, err)
	}
	return nil
}

// Sum 查询某个字段的和
func (repo *Repository) Sum(sum *finder.Sum) error {
	var db = repo.DB
	if sum.Debug {
		db.Debug()
	}
	if err := db.
		Model(sum.Model).
		Scopes(
			Wheres(sum.Wheres.Wheres...),
			Sum(sum.Col, sum.Val),
		).
		Error; err != nil {
		return err
	}
	return nil
}

func (repo *Repository) GetTableName(data any) string {
	return repo.DB.DB.NamingStrategy.TableName(reflect.TypeOf(data).Name())
}

func (repo *Repository) CreateTable(tableName string, data any) error {
	if len(tableName) > 0 {
		return repo.DB.Table(tableName).AutoMigrate(data)
	} else {
		return repo.DB.AutoMigrate(data)
	}
}

func (repo *Repository) TableHasExist(tableName string) bool {
	return repo.DB.DB.Migrator().HasTable(tableName)
}
