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
		conn, err = gorm.Open(postgres.New(postgres.Config{DSN: database.DSN, PreferSimpleProtocol: true}), &cfg)
	case conf.Sqlite:
		conn, err = gorm.Open(sqlite.Open(database.DSN), &cfg)
	default:
		conn, err = gorm.Open(mysql.Open(database.DSN), &cfg)
	}
	if err != nil {
		return nil, err
	}
	db, _ := conn.DB()
	if database.MaxOpenConns == 0 {
		database.MaxOpenConns = 20
	}
	if database.MaxIdleConns == 0 {
		database.MaxIdleConns = 10
	}
	db.SetMaxIdleConns(database.MaxIdleConns)
	db.SetMaxOpenConns(database.MaxOpenConns)
	db.SetConnMaxIdleTime(time.Minute * 10)
	db.SetConnMaxLifetime(time.Minute * 30)
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

type schemata struct {
	SchemaName string `gorm:"column:schema_name"`
}

/*
SELECT schema_name
FROM information_schema.schemata
-- 排除常见系统 Schema
WHERE schema_name NOT IN (

	                      'pg_catalog',    -- PG 核心系统表Schema
	                      'information_schema', -- SQL 标准系统视图Schema
	                      'pg_toast',      -- 大对象存储Schema
	                      'pg_temp_1',     -- 临时表Schema（会话级）
	                      'pg_toast_temp_1', -- 临时大对象Schema
	                     'public'
	)

ORDER BY schema_name;
*/
func GetAllSchema(db *gorm.DB) ([]string, error) {
	var s []schemata
	err := db.Table("information_schema.schemata").Select("schema_name").Where(` schema_name NOT IN (
                          'pg_catalog',    -- PG 核心系统表Schema
                          'information_schema', -- SQL 标准系统视图Schema
                          'pg_toast',      -- 大对象存储Schema
                          'pg_temp_1',     -- 临时表Schema（会话级）
                          'pg_toast_temp_1' -- 临时大对象Schema
    )`).Find(&s).Error
	if err != nil {
		return nil, err
	}
	var schemas []string
	for _, v := range s {
		schemas = append(schemas, v.SchemaName)
	}
	return schemas, nil
}
func SchemaTableAutoMigrate(ctx context.Context, schemas []string, dst ...interface{}) error {
	if len(schemas) == 0 {
		return errors.System.AddDetail("没有传入schema")
	}

	taskPool := make(chan struct{}, 20)
	var wg sync.WaitGroup
	var startTime = time.Now()
	for _, schema := range schemas {
		v := schema
		wg.Add(1)
		taskPool <- struct{}{}
		go func(schema string) {
			defer func() {
				<-taskPool
				wg.Done()
			}()
			db := GetSchemaTenantConn(SetTenantCode(ctx, schema))
			err := db.AutoMigrate(dst...)
			if err != nil {
				logx.WithContext(ctx).Errorf("SchemaTableAutoMigrate err:%v", err)
			} else {
				logx.WithContext(ctx).Infof("SchemaTableAutoMigrate finish schema:%v", schema)
			}
		}(v)
	}
	wg.Wait()
	logx.WithContext(ctx).Infof("SchemaTableAutoMigrate all finish use time:%v", time.Since(startTime))
	return nil
}

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
			cfg := gorm.Config{DisableForeignKeyConstraintWhenMigrating: true, DisableAutomaticPing: true, PrepareStmt: true, Logger: conn.Logger,
				NamingStrategy: schema.NamingStrategy{SingularTable: true, TablePrefix: cast.ToString(val) + "."}, ConnPool: conn.ConnPool}
			cc, err = gorm.Open(conn.Dialector, &cfg)
			if err != nil {
				conn.Error = errors.System.AddDetail("创建数据库连接失败").AddDetail(err)
				return conn
			}
			cc.Statement = conn.Statement
			cc.ConnPool = conn.ConnPool
			schemaMap.Store(tk, cc)
		} else {
			cc = c.(*gorm.DB).WithContext(ctx)
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
