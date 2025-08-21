package stores

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"github.com/glebarez/sqlite"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
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
		if err != nil {
			logx.Errorf("InitConn 失败 cfg:%v  err:%v", utils.Fmt(database), err)
			os.Exit(-1)
		}
		// 验证连接是否可用
		if sqlDB, err := commonConn.DB(); err == nil {
			if err := sqlDB.Ping(); err != nil {
				logx.Errorf("数据库连接ping失败: %v", err)
				os.Exit(-1)
			}
		}
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
		if err != nil {
			logx.Errorf("InitTsConn 失败: %v", err)
			os.Exit(-1)
		}
		// 验证时序数据库连接是否可用
		if sqlDB, err := tsConn.DB(); err == nil {
			if err := sqlDB.Ping(); err != nil {
				logx.Errorf("时序数据库连接ping失败: %v", err)
				os.Exit(-1)
			}
		}
		tsDBType = database.DBType
	})
	return
}

func GetConn(database conf.Database) (conn *gorm.DB, err error) {
	rlDBType = database.DBType
	cfg := gorm.Config{DisableForeignKeyConstraintWhenMigrating: true, PrepareStmt: true, Logger: NewLog(logger.Warn),
		NamingStrategy: schema.NamingStrategy{SingularTable: true}}
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
	dbCtxDebugKey      = "db.debug.type"
	dbCtxTenantCodeKey = "db.tenantCode"
)

func SetIsDebug(ctx context.Context, isDebug bool) context.Context {
	return context.WithValue(ctx, dbCtxDebugKey, isDebug)
}

func SetTenantCode(ctx context.Context, tenantCode string) context.Context {
	return context.WithValue(ctx, dbCtxTenantCodeKey, tenantCode)
}

func WithNoDebug[t any](ctx context.Context, f func(in any) t) t {
	ctx = SetIsDebug(ctx, false)
	return f(ctx)
}

// validateAndGetContext 验证输入参数并返回context，如果失败返回带错误的连接
func validateAndGetContext(in any, fallbackConn *gorm.DB) (context.Context, *gorm.DB) {
	if db, ok := in.(*gorm.DB); ok {
		return nil, db
	}

	ctx, ok := in.(context.Context)
	if !ok {
		// 如果传入的既不是*gorm.DB也不是context.Context，返回带错误的连接
		conn := fallbackConn.Session(&gorm.Session{})
		conn.Error = errors.Parameter.AddDetail("参数类型错误，需要*gorm.DB或context.Context")
		return nil, conn
	}

	return ctx, nil
}

// 获取租户连接  传入context或db连接 如果传入的是db连接则直接返回db
func GetTenantConn(in any) *gorm.DB {
	ctx, conn := validateAndGetContext(in, commonConn)
	if conn != nil {
		return conn // 返回错误连接或直接传入的DB连接
	}

	if val := ctx.Value(dbCtxDebugKey); val != nil && cast.ToBool(val) == false {
		return commonConn.WithContext(ctx)
	}
	return commonConn.WithContext(ctx).Debug()
}

// 获取租户连接 pg: schema级别隔离 mysql: database

var schemaMap sync.Map

func GetSchemaTenantConn(in any) *gorm.DB {
	ctx, conn := validateAndGetContext(in, commonConn)
	if conn != nil {
		return conn // 返回错误连接或直接传入的DB连接
	}
	conn = commonConn.WithContext(ctx)
	if val := ctx.Value(dbCtxTenantCodeKey); val != nil {
		tk := cast.ToString(val)
		c, ok := schemaMap.Load(tk)
		var cc *gorm.DB
		if !ok {
			var err error
			cfg := gorm.Config{DisableForeignKeyConstraintWhenMigrating: true, PrepareStmt: true, Logger: NewLog(logger.Warn),
				NamingStrategy: schema.NamingStrategy{SingularTable: true, TablePrefix: cast.ToString(val) + "."}}
			cc, err = gorm.Open(conn.Dialector, &cfg)
			if err != nil {
				conn.Error = errors.System.AddDetail("创建数据库连接失败").AddDetail(err)
				return conn
			}
			schemaMap.Store(tk, cc)
		} else {
			cc = c.(*gorm.DB)
		}
		if val := ctx.Value(dbCtxDebugKey); val != nil && cast.ToBool(val) == false {
			return cc
		}
		return cc.Debug()
	} else {
		conn.Error = errors.Permissions.AddDetail("没有传入租户号")
		return conn
	}

}

// 获取公共连接 传入context或db连接 如果传入的是db连接则直接返回db
func GetCommonConn(in any) *gorm.DB {
	ctx, conn := validateAndGetContext(in, commonConn)
	if conn != nil {
		return conn // 返回错误连接或直接传入的DB连接
	}

	if val := ctx.Value(dbCtxDebugKey); val != nil && cast.ToBool(val) == false {
		return commonConn.WithContext(ctx)
	}
	return commonConn.WithContext(ctx).Debug()
}

// 获取时序连接 传入context或db连接 如果传入的是db连接则直接返回db
func GetTsConn(in any) *gorm.DB {
	// 对于时序连接，需要特殊处理未初始化的情况
	if db, ok := in.(*gorm.DB); ok {
		return db
	}

	ctx, ok := in.(context.Context)
	if !ok {
		// 时序连接的特殊处理
		if tsConn == nil {
			conn := commonConn.Session(&gorm.Session{})
			conn.Error = errors.System.AddDetail("时序数据库连接未初始化")
			return conn
		}
		conn := tsConn.Session(&gorm.Session{})
		conn.Error = errors.Parameter.AddDetail("参数类型错误，需要*gorm.DB或context.Context")
		return conn
	}

	if tsConn == nil {
		conn := commonConn.WithContext(ctx)
		conn.Error = errors.System.AddDetail("时序数据库连接未初始化")
		return conn
	}
	if val := ctx.Value(dbCtxDebugKey); val != nil && cast.ToBool(val) == false {
		return tsConn.WithContext(ctx)
	}
	return tsConn.WithContext(ctx).Debug()
}

func SetAuthIncrement(conn *gorm.DB, table schema.Tabler) error {
	if conn == nil {
		return errors.Parameter.AddDetail("数据库连接不能为空")
	}
	if table == nil {
		return errors.Parameter.AddDetail("表结构不能为空")
	}

	var count int64
	err := conn.Model(table).Count(&count).Error
	if err != nil {
		return ErrFmt(err)
	}
	if count > 0 {
		return nil
	}

	// 使用新的会话避免影响原连接的状态
	session := conn.Session(&gorm.Session{})
	if err := session.Statement.Parse(table); err != nil {
		return ErrFmt(err)
	}

	err = session.Model(table).Create(table).Error
	if err != nil {
		return ErrFmt(err)
	}

	prioritizedPrimaryField := session.Statement.Schema.PrioritizedPrimaryField
	if prioritizedPrimaryField == nil {
		return nil
	}

	err = session.Model(table).Delete(table).Where(fmt.Sprintf("%s = ?", prioritizedPrimaryField.DBName), 10).Error
	return ErrFmt(err)
}

type PrefixClause struct {
	Prefix string
}

func (p PrefixClause) ModifyStatement(statement *gorm.Statement) {
	return
}

func (p PrefixClause) Name() string {
	return ""
}

func (p PrefixClause) Build(builder clause.Builder) {
	return
}

func (p PrefixClause) MergeClause(clause *clause.Clause) {
	return
}
