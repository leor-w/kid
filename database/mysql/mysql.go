package mysql

import (
	"context"
	"fmt"
	"github.com/leor-w/kid/config"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Conn struct {
	*gorm.DB
	options *Options
}

func (conn *Conn) Provide() interface{} {
	if !config.Exist("mysql") {
		panic("not found [mysql] in config")
	}
	return Default()
}

type Option func(*Options)

func New(opts ...Option) *Conn {
	options := Options{
		Port:    3306,
		MaxIdle: 20,
		MaxOpen: 100,
		MaxLife: 12,
	}
	for _, opt := range opts {
		opt(&options)
	}
	conn := &Conn{
		options: &options,
	}
	dns := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conn.options.User,
		conn.options.Password,
		conn.options.Host,
		conn.options.Port,
		conn.options.Db)
	db, err := gorm.Open(mysql.Open(dns), &gorm.Config{Logger: NewKidLogger(&logger.Config{
		SlowThreshold: 200 * time.Millisecond,
	})})
	if err != nil {
		panic(fmt.Sprintf("connect mysql failed: %s", err.Error()))
	}
	db.Logger.LogMode(logger.LogLevel(conn.options.LogLevel))

	sqlDb, err := db.DB()
	if err != nil {
		panic(fmt.Sprintf("mysql connect get database failed: %s", err.Error()))
	}
	sqlDb.SetMaxIdleConns(conn.options.MaxIdle)
	sqlDb.SetConnMaxLifetime(conn.options.MaxLife)
	sqlDb.SetMaxOpenConns(conn.options.MaxOpen)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err = sqlDb.PingContext(ctx); err != nil {
		panic(fmt.Sprintf("ping mysql failed: %s", err.Error()))
	}
	conn.DB = db
	return conn
}

func Default() *Conn {
	return New(
		WithHost(config.GetString("mysql.host")),
		WithPort(config.GetInt("mysql.port")),
		WithUser(config.GetString("mysql.user")),
		WithPassword(config.GetString("mysql.password")),
		WithDb(config.GetString("mysql.database")),
		WithMaxLife(config.GetDuration("mysql.maxLife")),
		WithMaxIdle(config.GetInt("mysql.maxIdle")),
		WithMaxOpen(config.GetInt("mysql.maxOpen")),
		WithLogLevel(config.GetInt("mysql.logLevel")),
	)
}
