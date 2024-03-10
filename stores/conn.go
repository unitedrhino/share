package stores

import (
	"context"
	"database/sql"
	"fmt"
	"gitee.com/i-Things/share/clients"
	"gitee.com/i-Things/share/conf"
	"github.com/dtm-labs/client/dtmgrpc"
	"github.com/glebarez/sqlite"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"sync"
	"time"
)

var (
	commonConn *gorm.DB
	once       sync.Once
	tenantConn sync.Map
	dbType     string //数据库类型
)

func InitConn(database conf.Database) {
	var err error
	once.Do(func() {
		commonConn, err = GetConn(database)
		logx.Must(err)
	})
	return
}

func GetConn(database conf.Database) (conn *gorm.DB, err error) {
	dbType = database.DBType
	switch database.DBType {
	case conf.Pgsql:
		conn, err = gorm.Open(postgres.Open(database.DSN), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true})
	case conf.Sqlite:
		conn, err = gorm.Open(sqlite.Open(database.DSN), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true})
	default:
		conn, err = gorm.Open(mysql.Open(database.DSN), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true})
	}
	if err != nil {
		return nil, err
	}
	db, _ := conn.DB()
	db.SetMaxIdleConns(200)
	db.SetMaxOpenConns(200)
	db.SetConnMaxIdleTime(time.Hour)
	db.SetConnMaxLifetime(time.Hour)
	return
}

func GetConnDB(database conf.Database) (db *sql.DB, err error) {
	switch database.DBType {
	case conf.Tdengine:
		td, err := clients.NewTDengine(database)
		if err != nil {
			return nil, err
		}
		return td.DB, nil
	default:
		conn, err := GetConn(database)
		if err != nil {
			return nil, err
		}
		return conn.DB()
	}
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

// 屏障分布式事务
func BarrierTransaction(ctx context.Context, fc func(tx *gorm.DB) error) error {
	conn := GetCommonConn(ctx)
	barrier, _ := dtmgrpc.BarrierFromGrpc(ctx)
	if barrier == nil { //如果没有开启分布式事务,则直接走普通事务即可
		return conn.Transaction(fc)
	}
	db, _ := conn.DB()
	return barrier.CallWithDB(db, func(tx *sql.Tx) error {
		gdb, _ := gorm.Open(mysql.New(mysql.Config{
			Conn: tx,
		}), &gorm.Config{})
		return fc(gdb)
	})
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
