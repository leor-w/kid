package mysql

import (
	"context"
	"fmt"
	"time"

	"github.com/lsjhtang/kid/config"
	"github.com/lsjhtang/kid/database/repos/join"
	"github.com/lsjhtang/kid/database/repos/where"
	"github.com/lsjhtang/kid/utils"

	"github.com/lsjhtang/kid/injector"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type MySQL struct {
	*gorm.DB
	options *Options
}

func (conn *MySQL) Provide(ctx context.Context) interface{} {
	var confName string
	name, ok := ctx.Value(injector.NameKey{}).(string)
	if ok && len(name) > 0 {
		confName = "." + name
	}
	confPrefix := fmt.Sprintf("mysql%s", confName)
	if !config.Exist(confPrefix) {
		panic(fmt.Sprintf("config.yaml file not found configuration item [%s]", confPrefix))
	}
	return New(
		WithHost(config.GetString(utils.GetConfigurationItem(confPrefix, "host"))),
		WithPort(config.GetInt(utils.GetConfigurationItem(confPrefix, "port"))),
		WithUser(config.GetString(utils.GetConfigurationItem(confPrefix, "user"))),
		WithPassword(config.GetString(utils.GetConfigurationItem(confPrefix, "password"))),
		WithDb(config.GetString(utils.GetConfigurationItem(confPrefix, "database"))),
		WithMaxLife(config.GetDuration(utils.GetConfigurationItem(confPrefix, "maxLife"))),
		WithMaxIdle(config.GetInt(utils.GetConfigurationItem(confPrefix, "maxIdle"))),
		WithMaxOpen(config.GetInt(utils.GetConfigurationItem(confPrefix, "maxOpen"))),
		WithLogLevel(config.GetInt(utils.GetConfigurationItem(confPrefix, "logLevel"))),
	)
}

type Option func(*Options)

func New(opts ...Option) *MySQL {
	options := Options{
		Port:    3306,
		MaxIdle: 20,
		MaxOpen: 100,
		MaxLife: 12,
	}
	for _, opt := range opts {
		opt(&options)
	}
	conn := &MySQL{
		options: &options,
	}
	dns := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conn.options.User,
		conn.options.Password,
		conn.options.Host,
		conn.options.Port,
		conn.options.Db)
	db, err := gorm.Open(mysql.Open(dns), &gorm.Config{
		Logger: NewKidLogger(&logger.Config{
			SlowThreshold: 200 * time.Millisecond,
			Colorful:      true,
			LogLevel:      logger.LogLevel(conn.options.LogLevel),
		}),
		DisableForeignKeyConstraintWhenMigrating: options.CloseFKCheck,
	})
	if err != nil {
		panic(fmt.Sprintf("开启MySQL连接错误: %s", err.Error()))
	}

	sqlDb, err := db.DB()
	if err != nil {
		panic(fmt.Sprintf("创建MySQL连接对象错误: %s", err.Error()))
	}
	sqlDb.SetMaxIdleConns(conn.options.MaxIdle)
	sqlDb.SetConnMaxLifetime(conn.options.MaxLife)
	sqlDb.SetMaxOpenConns(conn.options.MaxOpen)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err = sqlDb.PingContext(ctx); err != nil {
		panic(fmt.Sprintf("测试MySQL连接状态错误: %s", err.Error()))
	}
	conn.DB = db
	return conn
}

func Paginate(pageNum, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if pageNum <= 0 {
			pageNum = 1
		}
		if pageSize <= 0 {
			pageSize = 10
		}
		if pageSize >= 100 {
			pageSize = 100
		}
		offset := (pageNum - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

type WhereWrapper struct {
	Column    string
	Condition string
	Value     interface{}
}

func NewWhereWrapper(col string, condition string, val interface{}) *WhereWrapper {
	return &WhereWrapper{
		Column:    col,
		Condition: condition,
		Value:     val,
	}
}

func Wrappers(conditions ...*WhereWrapper) []*WhereWrapper {
	return conditions
}

const (
	ConditionTypeLte = "<="
	ConditionTypeLt  = "<"
	ConditionTypeGte = ">="
	ConditionTypeGt  = ">"
	ConditionTypeEq  = "="
	ConditionTypeIn  = "in"
)

func Wheres(wheres ...*where.Where) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(wheres) <= 0 {
			return db
		}
		for _, w := range wheres {
			format, vals := w.Build()
			if w.Or {
				db.Or(format, vals...)
				continue
			}
			db.Where(format, vals...)
		}
		return db
	}
}

func Joins(joins ...join.Join) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for _, v := range joins {
			db.Joins(v.Build())
		}
		return db
	}
}

func Sum(column string, sum interface{}) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Select(fmt.Sprintf("IFNULL(SUM(%s), 0) AS sum", column)).Scan(&sum)
	}
}

func Where(wheres []*WhereWrapper) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(wheres) <= 0 {
			return db
		}
		for _, wrapper := range wheres {
			if wrapper.Condition == ConditionTypeIn {
				db.Where(fmt.Sprintf("%s IN(?)", wrapper.Column), wrapper.Value)
				continue
			}
			db.Where(fmt.Sprintf("%s %s ?", wrapper.Column, wrapper.Condition), wrapper.Value)
		}
		return db
	}
}

func SearchScope(wheres where.Wheres) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(wheres.Wheres) <= 0 {
			return db
		}
		for _, w := range wheres.Wheres {
			var notNil bool
			switch vt := w.V1.(type) {
			case string:
				notNil = vt != ""
			case uint8, int8, uint, int, uint32, int32, uint16, int16, uint64, int64, float32, float64:
				notNil = vt != 0
			}
			if notNil {
				db.Or(fmt.Sprintf("%s LIKE '%v'", w.Column, w.V1))
			}
		}
		return db
	}
}

// Search 搜索 多个值 当有值
func Search(wheres []*WhereWrapper) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(wheres) <= 0 {
			return db
		}
		for _, wrapper := range wheres {
			var notNil bool
			switch vt := wrapper.Value.(type) {
			case string:
				notNil = vt != ""
			case uint8, int8, uint, int, uint32, int32, uint16, int16, uint64, int64, float32, float64:
				notNil = vt != 0
			}
			if notNil {
				db.Or(fmt.Sprintf("%s LIKE '%%%v%%'", wrapper.Column, wrapper.Value))
			}
		}
		return db
	}
}

func FilterDeleted() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("deleted = 0")
	}
}

type AutoMigrate interface {
	Models() []interface{}
}
