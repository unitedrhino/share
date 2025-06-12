package stores

import (
	"context"
	"fmt"
	"gitee.com/unitedrhino/share/conf"
	"github.com/glebarez/sqlite"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"sync"
	"time"
)

type DB = gorm.DB

var (
	commonConn *gorm.DB
	tsConn     *gorm.DB
	once       sync.Once
	tsOnce     sync.Once
	tenantConn sync.Map
	rlDBType   string                 //关系型数据库类型
	tsDBType   string = conf.Tdengine //时序数据库类型
)

func GetTsDBType() string {
	return tsDBType
}
func GetDBType() string {
	return rlDBType
}

func InitConn(database conf.Database) {
	var err error
	once.Do(func() {
		commonConn, err = GetConn(database)
		logx.Must(err)
	})
	return
}
func InitTsConn(database conf.TSDB) {
	var err error
	tsOnce.Do(func() {
		tsConn, err = GetConn(conf.Database{
			DBType:      database.DBType,
			IsInitTable: true,
			DSN:         database.DSN,
		})
		logx.Must(err)
		tsDBType = database.DBType
	})
	return
}

func GetConn(database conf.Database) (conn *gorm.DB, err error) {
	rlDBType = database.DBType
	cfg := gorm.Config{DisableForeignKeyConstraintWhenMigrating: true, PrepareStmt: true, Logger: NewLog(logger.Warn)}
	switch database.DBType {
	case conf.Pgsql:
		conn, err = gorm.Open(postgres.Open(database.DSN), &cfg)
	case conf.Sqlite:
		conn, err = gorm.Open(sqlite.Open(database.DSN), &cfg)
	default:
		conn, err = gorm.Open(mysql.Open(database.DSN), &cfg)
	}
	if err != nil {
		return nil, err
	}
	db, _ := conn.DB()
	db.SetMaxIdleConns(20)
	db.SetMaxOpenConns(20)
	db.SetConnMaxIdleTime(time.Hour)
	db.SetConnMaxLifetime(time.Hour)
	return
}

const (
	dbCtxDebugKey = "db.debug.type"
)

func SetIsDebug(ctx context.Context, isDebug bool) context.Context {
	return context.WithValue(ctx, dbCtxDebugKey, isDebug)
}

func WithNoDebug[t any](ctx context.Context, f func(in any) t) t {
	ctx = SetIsDebug(ctx, false)
	return f(ctx)
}

// 获取租户连接  传入context或db连接 如果传入的是db连接则直接返回db
func GetTenantConn(in any) *gorm.DB {
	if db, ok := in.(*gorm.DB); ok {
		return db
	}
	ctx := in.(context.Context)
	if val := ctx.Value(dbCtxDebugKey); val != nil && cast.ToBool(val) == false { //不打印日志
		return commonConn.WithContext(ctx)
	}
	return commonConn.WithContext(ctx).Debug()
}

// 获取公共连接 传入context或db连接 如果传入的是db连接则直接返回db
func GetCommonConn(in any) *gorm.DB {
	if db, ok := in.(*gorm.DB); ok {
		return db
	}
	ctx := in.(context.Context)
	if val := ctx.Value(dbCtxDebugKey); val != nil && cast.ToBool(val) == false { //不打印日志
		return commonConn.WithContext(ctx)
	}
	return commonConn.WithContext(ctx).Debug()
}

// 获取时序连接 传入context或db连接 如果传入的是db连接则直接返回db
func GetTsConn(in any) *gorm.DB {
	if db, ok := in.(*gorm.DB); ok {
		return db
	}
	ctx := in.(context.Context)
	if val := ctx.Value(dbCtxDebugKey); val != nil && cast.ToBool(val) == false { //不打印日志
		return tsConn.WithContext(ctx)
	}
	return tsConn.WithContext(ctx).Debug()
}

func SetAuthIncrement(conn *gorm.DB, table schema.Tabler) error {
	var count int64
	err := conn.Model(table).Count(&count).Error
	if err != nil {
		return ErrFmt(err)
	}
	if count > 0 {
		return nil
	}
	if err := conn.Statement.Parse(table); err != nil {
		return err
	}
	err = conn.Model(table).Create(table).Error
	if err != nil {
		return err
	}
	prioritizedPrimaryField := conn.Statement.Schema.PrioritizedPrimaryField
	if prioritizedPrimaryField == nil {
		return nil
	}
	err = conn.Model(table).Delete(table).Where(fmt.Sprintf("%s = ?", prioritizedPrimaryField.DBName), 10).Error
	return err
}
